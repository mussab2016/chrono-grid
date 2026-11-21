import {Component, inject, OnInit, signal, computed} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormBuilder, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatIconModule} from '@angular/material/icon';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatProgressBarModule} from '@angular/material/progress-bar';

import {ShellService} from '@common/shell-service/shell.service';
import {TranslatePipe} from '../../pipes/unix-transtale.pipe';
import {AcademicYear} from '../../../api/academic-year';
import {AcademicTerm} from '../../../api/academic-term';
import {AcademicYearService} from '../academic-year.service';
import {AcademicTermService} from '../academic-term.service';
import {StudentGroup} from '../../../api/student-group';
import {AcademicLevel} from '../../../api/academic-level';
import {StudentGroupService} from '../student-group.service';
import {AcademicLevelService} from '../../academic-levels/academic-level.service';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
    selector: 'app-academic-calendar',
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatIconModule,
        MatTooltipModule,
        MatProgressBarModule,
        TranslatePipe
    ],
    templateUrl: './academic-calendar.component.html',
    styleUrl: './academic-calendar.component.scss'
})
export class AcademicCalendarComponent implements OnInit {
    protected router = inject(Router);
    protected route = inject(ActivatedRoute);
    private yearSrv = inject(AcademicYearService);
    private termSrv = inject(AcademicTermService);
    private majorsSrv = inject(AcademicLevelService);
    private groupsSrv = inject(StudentGroupService);
    private shellSrv = inject(ShellService);
    private fb = inject(FormBuilder);

    // State Signals
    years = signal<AcademicYear[]>([]);
    selectedYear = signal<AcademicYear | null>(null);
    selectedTerm = signal<AcademicTerm | null>(null);
    selectedMajor = signal<AcademicLevel | null>(null);
    selectedGroup = signal<StudentGroup | null>(null);
    terms = signal<AcademicTerm[]>([]);
    loadingYears = signal(false);
    loadingTerms = signal(false);
    loadingGroups = signal(false);
    loadingMajors = signal(false);
    groups = signal<StudentGroup[]>([]);
    majors = signal<AcademicLevel[]>([]);

    // Inline Quick-Add / Edit Controls
    newYearCtrl = this.fb.control('', [Validators.required, Validators.maxLength(20)]);
    editingTermId = signal<number | null>(null);
    // editingGroupId = signal<number | null>(null);

    termForm = this.fb.group({
        semester: [1, [Validators.required, Validators.min(1), Validators.max(12)]],
        label: ['', [Validators.required, Validators.maxLength(20)]]
    });

    newGroupCtrl = this.fb.control('', [Validators.required, Validators.maxLength(20)])

    // auto suggesting if the terms based on the existing terms
    nextSuggestedSemester = computed(() => {
        const list = this.terms();
        if (list.length === 0) return 1;
        return Math.max(...list.map(t => t.semester)) + 1;
    });

    ngOnInit(): void {
        this.shellSrv.setTitle('ACADEMIC_CALENDAR');
        this.shellSrv.setGoBack(this.exit.bind(this));
        this.loadYears();
        this.loadMajors();
    }

    loadYears(selectId?: number): void {
        this.loadingYears.set(true);
        this.yearSrv.getAll('', '-id').subscribe({
            next: (res) => {
                const list = res.data ?? [];
                this.years.set(list);
                this.loadingYears.set(false);

                if (list.length > 0) {
                    const target = selectId
                        ? list.find(y => y.id === selectId) ?? list[0]
                        : this.selectedYear()
                            ? list.find(y => y.id === this.selectedYear()!.id) ?? list[0]
                            : list[0];
                    this.selectYear(target);
                } else {
                    this.selectedYear.set(null);
                    this.terms.set([]);
                }
            },
            error: () => this.loadingYears.set(false)
        });
    }

    selectYear(year: AcademicYear): void {
        this.selectedTerm.set(null)
        this.selectedMajor.set(null)
        this.selectedYear.set(year);
        this.selectedGroup.set(null)
        this.cancelEditTerm();
        this.loadTerms(year.id);
        this.loadGroups(year.id); // Fetch all groups for this academic year once
    }

    loadMajors(): void {
        this.loadingMajors.set(true);
        this.majorsSrv.getAll('', '-id').subscribe({
            next: (res) => {
                const list = res.data ?? [];
                this.majors.set(list);
                this.loadingMajors.set(false);
            },
            error: () => this.loadingMajors.set(false)
        });
    }

    selectTerm(term: AcademicTerm): void {
        this.selectedMajor.set(null)
        this.selectedTerm.set(term);
        this.cancelEditTerm();
        // this.loadTerms(year.id); //this should be changed to upload session offerings
    }

    loadTerms(yearId: number): void {
        this.loadingTerms.set(true);
        this.termSrv.getByYearId(yearId).subscribe({
            next: (response) => {
                this.terms.set(response.data ?? []);
                this.resetTermFormToNext();
                this.loadingTerms.set(false);
            },
            error: (err) => {
                this.loadingTerms.set(false)
                throw err
            }
        });
    }

    loadGroups(yearId: number): void {
        this.loadingGroups.set(true);
        this.groupsSrv.getByYearId(yearId).subscribe({
            next: (response) => {
                this.groups.set(response.data ?? []);
                this.loadingGroups.set(false);
            },
            error: (err) => {
                this.loadingGroups.set(false)
                throw err
            }
        });
    }

    // --- years managements ---
    addYear(): void {
        const label = this.newYearCtrl.value?.trim();
        if (!label) return;

        this.yearSrv.save({id: 0, yearLabel: label}).subscribe({
            next: (id) => {
                this.newYearCtrl.reset();
                this.loadYears(id);
            }
        });
    }

    deleteYear(year: AcademicYear, event: Event): void {
        event.stopPropagation();
        this.yearSrv.remove(year.id).subscribe({
            next: () => {
                if (this.selectedYear()?.id === year.id) {
                    this.selectedYear.set(null);
                }
                this.loadYears();
            }
        });
    }

    // --- terms managements ---
    resetTermFormToNext(): void {
        const nextSem = this.nextSuggestedSemester();
        this.editingTermId.set(null);
        this.termForm.reset({
            semester: nextSem,
            label: `S${nextSem}`
        });
    }

    editTerm(term: AcademicTerm): void {
        this.editingTermId.set(term.id);
        this.termForm.patchValue({
            semester: term.semester,
            label: term.label
        });
    }

    cancelEditTerm(): void {
        this.resetTermFormToNext();
    }

    saveTerm(): void {
        const currentYear = this.selectedYear();
        if (!currentYear || this.termForm.invalid) return;

        const payload: AcademicTerm = {
            id: this.editingTermId() ?? 0,
            academicYearId: currentYear.id,
            semester: Number(this.termForm.value.semester),
            label: this.termForm.value.label!.trim()
        };

        this.termSrv.save(payload).subscribe({
            next: () => {
                this.loadTerms(currentYear.id);
            }
        });
    }

    deleteTerm(term: AcademicTerm): void {
        const currentYear = this.selectedYear();
        if (!currentYear) return;

        this.termSrv.remove(term.id).subscribe({
            next: () => this.loadTerms(currentYear.id)
        });
    }

    // --- StudentGroups managements ---
    addGroup(): void {
        if (!this.selectedMajor() || !this.selectedYear()) {
            return;
        }
        const yearId = this.selectedYear()!.id;
        const code = this.newGroupCtrl.value?.trim();
        if (!code) return;

        const studentGroup: StudentGroup = {
            id: 0,
            academicYearID: this.selectedYear()!.id,
            levelId: this.selectedMajor()!.id,
            code: code,
        }
        this.groupsSrv.save(studentGroup).subscribe({
            next: () => {
                this.newGroupCtrl.reset();
                this.loadGroups(yearId);
            }
        });
    }

    deleteGroup(group: StudentGroup, event: Event): void {
        event.stopPropagation();
        if (!this.selectedMajor() || !this.selectedYear()) {
            return;
        }
        const yearId = this.selectedYear()!.id;
        this.groupsSrv.remove(group.id).subscribe({
            next: () => {
                if (this.selectedGroup()?.id === group.id) {
                    this.selectedGroup.set(null); // Todo remove stop stopPropagation. when want to edit or delete he should select
                }
                this.loadGroups(yearId);
            }
        });
    }

    // Filtered groups for the selected major
    filteredGroups = computed(() => {
        const major = this.selectedMajor();
        const all = this.groups();
        if (!major) return [];
        const l = all.filter(g => g.levelId === major.id);
        console.log(l.length)
        return all.filter(g => g.levelId === major.id);
    });

    // Helper map to count groups per major
    majorGroupCounts = computed(() => {
        const counts = new Map<number, number>();
        for (const g of this.groups()) {
            counts.set(g.levelId, (counts.get(g.levelId) ?? 0) + 1);
        }
        return counts;
    });

    selectMajor(major: AcademicLevel): void {
        this.selectedMajor.set(major);
        this.selectedGroup.set(null);
        this.filteredGroups()
        this.newGroupCtrl.reset();
    }

    manageOffering(): void {
        const term = this.selectedTerm();
        if (!term) return;
        this.router.navigate([term.id], { relativeTo: this.route });
    }

    exit(): void {
        this.router.navigate(['..'], { relativeTo: this.route }).then();
    }
}
