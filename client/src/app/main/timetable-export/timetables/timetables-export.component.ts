import {Component, OnInit, computed, inject, signal, viewChild} from '@angular/core';
import {FormControl, FormGroup, ReactiveFormsModule} from '@angular/forms';
import {toSignal} from '@angular/core/rxjs-interop';
import {forkJoin} from 'rxjs';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatSelectModule} from '@angular/material/select';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';

import {TimetableSlot} from '../../../api/timetable';
import {EntityOption, TimetableViewType} from '../model/timetable.model';
import {TimetableComponent} from '../model/timetable.component';
import {TimetableService} from '../../timetable/timetable.service';
import {AcademicTermService} from '../../academic_years/academic-term.service';
import {AcademicTerm} from '../../../api/academic-term';
import {ProfessorService} from '../../professors/professor.service';
import {Professor} from '../../../api/professor';
import {RoomsService} from '../../rooms/rooms.service';
import {AcademicLevelService} from '../../academic-levels/academic-level.service';
import {StudentGroupService} from '../../academic_years/student-group.service';
import {DefaultAutoCompleteModel} from '@common/ui-auto-complete/auto-complete-model';
import {UiAutoCompleteComponent} from '@common/ui-auto-complete/auto-complete.component';
import {ActivatedRoute, Router} from '@angular/router';
import {ShellService} from '@common/shell-service/shell.service';
import {StudentGroup} from '../../../api/student-group';
import {AcademicLevel} from '../../../api/academic-level';
import {Room} from '../../../api/room';

@Component({
    selector: 'timetables-export-component',
    standalone: true,
    templateUrl: './timetables-export.component.html',
    styleUrl: './timetables-export.component.scss',
    imports: [
        ReactiveFormsModule,
        TimetableComponent,
        MatFormFieldModule,
        MatSelectModule,
        MatButtonModule,
        MatIconModule,
        UiAutoCompleteComponent,
    ]
})
export class TimetablesExportComponent implements OnInit {
    private readonly route = inject(ActivatedRoute);
    private readonly router= inject(Router);
    private readonly timetableSrv = inject(TimetableService);
    private readonly termsSrv = inject(AcademicTermService);
    private readonly roomSrv = inject(RoomsService);
    private readonly levelsSrv = inject(AcademicLevelService);
    private readonly groupsSrv = inject(StudentGroupService);
    private readonly professorSrv = inject(ProfessorService);
    private readonly shellSrv = inject(ShellService);

    // Server data signals (preserving exact backend order)
    readonly timetable = signal<TimetableSlot[]>([]);
    readonly termsList = signal<AcademicTerm[]>([]);
    readonly roomsList = signal<Room[]>([]);
    readonly academicLevelsList = signal<AcademicLevel[]>([]);
    readonly studentGroupList = signal<StudentGroup[]>([]);
    readonly professorsList = signal<Professor[]>([]);
    readonly isLoading = signal<boolean>(true);

    // Active view mode
    readonly selectedType = signal<TimetableViewType>('AcademicLevels');

    // Unified reactive filter form
    readonly filterForm = new FormGroup({
        termId: new FormControl<number | null>(null),
        levelId: new FormControl<number | null>(null),
        groupNumber: new FormControl<number | null>(null),
        roomId: new FormControl<number | null>(null),
        professorId: new FormControl<number | null>(null)
    });

    // Ready-made AutoComplete model for Professors
    readonly professorAutoComplete = new DefaultAutoCompleteModel<Professor>(
        this.professorSrv,
        'fullName',
        (p) => p.fullName,
        'id'
    );

    // Reactive signals synced directly with FormControls
    readonly selectedLevelId = toSignal(this.filterForm.controls.levelId.valueChanges, {initialValue: null});
    readonly selectedGroup = toSignal(this.filterForm.controls.groupNumber.valueChanges, {initialValue: null});
    readonly selectedRoomId = toSignal(this.filterForm.controls.roomId.valueChanges, {initialValue: null});
    readonly selectedProfessorId = toSignal(this.filterForm.controls.professorId.valueChanges, {initialValue: null});

    // Single computed targetId routed by active view type
    readonly activeTargetId = computed<number | null>(() => {
        const type = this.selectedType();
        if (type === 'Class') return this.selectedRoomId();
        if (type === 'Professor') return this.selectedProfessorId();
        return this.selectedLevelId();
    });

    // Groups belonging to the currently selected academic level
    readonly filteredGroups = computed<StudentGroup[]>(() => {
        const levelId = this.selectedLevelId();
        if (levelId === null) return [];
        return this.studentGroupList().filter(g => g.levelId === levelId);
    });

    // Capture the child component reference safely across @if/@else blocks
    readonly timetableRef = viewChild(TimetableComponent);

    ngOnInit(): void {
        this.shellSrv.setGoBack(this.exit.bind(this));
        this.shellSrv.setTitle('Timetables');
        this.loadInitialData();

        // Fetch timetable automatically when termId control changes
        this.filterForm.controls.termId.valueChanges.subscribe((termId) => {
            if (termId) {
                this.fetchTimetable(termId);
            }
        });

        // Reset group selection whenever academic level changes
        this.filterForm.controls.levelId.valueChanges.subscribe(() => {
            this.filterForm.controls.groupNumber.setValue(null);
        });
    }

    private loadInitialData(): void {
        this.isLoading.set(true);

        forkJoin({
            terms: this.termsSrv.getAll(),
            rooms: this.roomSrv.getAll(),
            levels: this.levelsSrv.getAll(),
            groups: this.groupsSrv.getAll(),
            professors: this.professorSrv.getAll()
        }).subscribe({
            next: ({terms, rooms, levels, groups, professors}) => {
                const termList = terms.data ?? [];
                this.termsList.set(termList);

                this.roomsList.set(rooms.data);

                this.academicLevelsList.set(levels.data);


                this.studentGroupList.set(groups.data);

                this.professorsList.set(professors.data);

                const firstTermId = termList[0]?.id ?? null;
                if (firstTermId !== null) {
                    this.filterForm.controls.termId.setValue(firstTermId);
                } else {
                    this.isLoading.set(false);
                }
            },
            error: () => {
                this.isLoading.set(false);
            }
        });
    }

    fetchTimetable(termId: number): void {
        this.isLoading.set(true);

        this.timetableSrv.getByTerm(termId).subscribe({
            next: (slots) => {
                this.timetable.set(slots ?? []);
                this.filterForm.controls.groupNumber.setValue(null);
                this.isLoading.set(false);
            },
            error: () => {
                this.isLoading.set(false);
            }
        });
    }

    setViewType(type: TimetableViewType): void {
        this.selectedType.set(type);
        this.filterForm.patchValue({
            levelId: null,
            groupNumber: null,
            roomId: null,
            professorId: null
        });
    }

    listToEntityOption(list: readonly (Professor | Room)[]): EntityOption[] {
        if (list.length === 0) {
            return [];
        }

        return list.map((item): EntityOption => ({
            id: item.id,
            label: 'fullName' in item ? item.fullName : item.code
        }));
    }

    exportExcel(): void {
        this.timetableRef()?.exportExcel();
    }

    printPage(): void {
        window.print();
    }

    public exit(): void {
        this.router.navigate(['..'], { relativeTo: this.route }).then();
    }
}
