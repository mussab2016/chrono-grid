import {Component, computed, inject, OnInit, signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {CommonModule} from '@angular/common';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatSelectModule} from '@angular/material/select';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {TimetableService} from '../timetable.service';
import {ProfessorService} from '../../professors/professor.service';
import {Professor} from '../../../api/professor';
import {ProfPreferenceInput, TimetableSlot} from '../../../api/timetable';
import {TranslatePipe} from '../../pipes/unix-transtale.pipe';
import {UiAutoCompleteComponent} from '@common/ui-auto-complete/auto-complete.component';
import {DefaultAutoCompleteModel} from '@common/ui-auto-complete/auto-complete-model';

@Component({
    selector: 'app-timetable-generator',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        RouterLink,
        MatButtonModule,
        MatSelectModule,
        MatFormFieldModule,
        MatProgressBarModule,
        MatSnackBarModule,
        TranslatePipe,
        ReactiveFormsModule,
        UiAutoCompleteComponent
    ],
    templateUrl: './timetable-generator.component.html',
    styleUrl: './timetable-generator.component.scss'
})
export class TimetableGeneratorComponent implements OnInit {
    private readonly route = inject(ActivatedRoute);
    private readonly timetableSrv = inject(TimetableService);
    private readonly professorSrv = inject(ProfessorService);
    private readonly snackBar = inject(MatSnackBar);

    termId = signal<number>(0);
    // professors = signal<Professor[]>([]);

    daysCount = signal<number>(5);
    slotsPerDay = signal<number>(6);
    generating = signal<boolean>(false);
    slots = signal<TimetableSlot[]>([]);

    // Temporary UX preferences map: key = "profId:day:slot" -> score (-1 | 10)
    selectedProfId = signal<number | null>(null);
    prefMap = signal<Record<string, number>>({});

    // View filter: '' (ALL) | group string
    filterGroup = signal<number>(-1);

    daysArray = computed(() => Array.from({length: this.daysCount()}, (_, i) => i + 1));
    slotsArray = computed(() => Array.from({length: this.slotsPerDay()}, (_, i) => i + 1));

    hasPreferences = computed(() => Object.keys(this.prefMap()).length > 0);

    availableGroups = computed(() => {
        const set = new Set<number>();
        for (const s of this.slots()) {
            for (const g of s.groups ?? []) {
                set.add(g);
            }
        }
        return Array.from(set).sort();
    });

    filteredSlots = computed(() => {
        const g = this.filterGroup();
        if (!g) return this.slots();
        return this.slots().filter(s => (s.groups ?? []).includes(g));
    });

    professorAutoComplete = new DefaultAutoCompleteModel<Professor>(
        this.professorSrv,
        'fullName',
        (p) => p.fullName,
        'id'
    );

    ngOnInit(): void {
        // Resolve termId from current or parent route snapshot
        const rawId = this.route.snapshot.paramMap.get('termId')
            ?? this.route.parent?.snapshot.paramMap.get('termId')
            ?? '0';
        const id = Number(rawId);
        this.termId.set(id);

        // this.professorSrv.getAll().subscribe({
        //     next: (list) => this.professors.set(list.data ?? [])
        // });

        if (id > 0) {
            this.loadExisting();
        }
    }

    // Fetch existing timetable slots for the active academic term
    loadExisting(): void {
        this.timetableSrv.getByTerm(this.termId()).subscribe({
            next: (data) => this.slots.set(data ?? [])
        });
    }

    // Cycle cell state on click: 0 (Neutral) -> 10 (Preferred) -> -1 (Forbidden) -> 0
    cyclePref(day: number, slot: number): void {
        const pId = this.selectedProfId();
        if (!pId) return;

        const key = `${pId}:${day}:${slot}`;
        const current = this.prefMap()[key] ?? 0;
        const next = current === 0 ? 10 : current === 10 ? -1 : 0;

        const updated = {...this.prefMap()};
        if (next === 0) {
            delete updated[key];
        } else {
            updated[key] = next;
        }
        this.prefMap.set(updated);
    }

    // Toggle entire day forbidden for selected professor
    toggleWholeDayForbidden(day: number): void {
        const pId = this.selectedProfId();
        if (!pId) return;

        const key = `${pId}:${day}:0`;
        const updated = {...this.prefMap()};
        if (updated[key] === -1) {
            delete updated[key];
        } else {
            updated[key] = -1;
        }
        this.prefMap.set(updated);
    }

    clearPreferences(): void {
        this.prefMap.set({});
    }

    getCellState(day: number, slot: number): number {
        const pId = this.selectedProfId();
        if (!pId) return 0;
        if (this.prefMap()[`${pId}:${day}:0`] === -1) return -1;
        return this.prefMap()[`${pId}:${day}:${slot}`] ?? 0;
    }

    getSlotsAt(day: number, slot: number): TimetableSlot[] {
        return this.filteredSlots().filter(s => s.dayNumber === day && s.slotNumber === slot);
    }

    // Trigger AI solver generation and update matrix slots
    generateNow(): void {
        const preferences: ProfPreferenceInput[] = Object.entries(this.prefMap()).map(([key, score]) => {
            const [professorId, dayNumber, slotNumber] = key.split(':').map(Number);
            return {professorId, dayNumber, slotNumber, score};
        });

        this.generating.set(true);
        this.timetableSrv.generate({
            termId: this.termId(),
            daysCount: this.daysCount(),
            slotsPerDay: this.slotsPerDay(),
            preferences
        }).subscribe({
            next: (res) => {
                this.generating.set(false);
                if (res.status === 'INFEASIBLE') {
                    this.snackBar.open('No conflict-free solution found. Relax some constraints.', 'OK', {duration: 4500});
                    return;
                }
                this.slots.set(res.slots ?? []);
                this.snackBar.open(`Timetable generated (${res.status})`, 'OK', {duration: 3000});
            },
            error: (err) => {
                this.generating.set(false);
                this.snackBar.open('Failed to generate timetable', 'OK', {duration: 4000});
                throw err;
            }
        });
    }
}
