import {Component, OnInit, inject, signal, computed} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ActivatedRoute, Router} from '@angular/router';
import {FormBuilder, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatProgressBarModule} from '@angular/material/progress-bar';

import {ShellService} from '@common/shell-service/shell.service';
import {TranslatePipe} from '../../pipes/unix-transtale.pipe';
import {StudentGroupService} from '../../academic_years/student-group.service';
import {CourseOfferingService} from '../course-offering.service';
import {AcademicTermService} from '../../academic_years/academic-term.service';
import {AcademicLevelService} from '../../academic-levels/academic-level.service';
import {SubjectService} from '../../subjects/subject.service';
import {AcademicLevel} from '../../../api/academic-level';
import {StudentGroup} from '../../../api/student-group';
import {CourseOfferingBatchRequest, CourseOfferingBlock} from '../../../api/course-offering';
import {AcademicTerm} from '../../../api/academic-term';
import {ProfessorService} from '../../professors/professor.service';
import {RoomTypeService} from '../../roomTypes/room-type.service';
import {TranslateDirective} from '@common/translate-directive/translate.directive';
import {UiAutoCompleteComponent} from '@common/ui-auto-complete/auto-complete.component';
import {DefaultAutoCompleteModel} from '@common/ui-auto-complete/auto-complete-model';
import {Professor} from '../../../api/professor';
import {Subject} from '../../../api/subject';

export interface SubjectColumn {
    subjectId: number;
    subjectCode: string;
    subjectName: string;
    blocks: CourseOfferingBlock[];
}

@Component({
    selector: 'app-course-offerings',
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatAutocompleteModule,
        MatSlideToggleModule,
        MatTooltipModule,
        MatProgressBarModule,
        TranslatePipe,
        TranslateDirective,
        UiAutoCompleteComponent
    ],
    templateUrl: './course-offerings.component.html',
    styleUrl: './course-offerings.component.scss'
})
export class CourseOfferingsComponent implements OnInit {
    private router = inject(Router);
    private route = inject(ActivatedRoute);
    private shellSrv = inject(ShellService);
    private fb = inject(FormBuilder);

    private offeringSrv = inject(CourseOfferingService);
    private termSrv = inject(AcademicTermService);
    private groupSrv = inject(StudentGroupService);
    private majorSrv = inject(AcademicLevelService);
    private subjectSrv = inject(SubjectService);
    private profSrv = inject(ProfessorService);
    private roomTypeSrv = inject(RoomTypeService);

    // Route & Context Signals
    termId = signal<number>(0);
    term = signal<AcademicTerm | null>(null);

    // Master Data Signals
    blocks = signal<CourseOfferingBlock[]>([]);
    // subjects = signal<any[]>([]);
    // professors = signal<any[]>([]);
    roomTypes = signal<any[]>([]);
    majors = signal<AcademicLevel[]>([]);
    groups = signal<StudentGroup[]>([]);

    // UI Filter & Builder Signals
    loading = signal(false);
    selectedMajorFilter = signal<number | null>(null);
    editingSessionUuid = signal<string | null>(null);
    selectedGroupIds = signal<number[]>([]);
    isBuilderOpen = signal(false);

    // AutoComplete Models (On-demand search + built-in cache)
    subjectAutoComplete = new DefaultAutoCompleteModel<Subject>(
        this.subjectSrv,
        'name',
        (s) => `${s.code} - ${s.name}`,
        'id'
    );

    professorAutoComplete = new DefaultAutoCompleteModel<Professor>(
        this.profSrv,
        'fullName',
        (p) => p.fullName,
        'id'
    );

    // Square Session Builder Form
    builderForm = this.fb.group({
        subjectId: [null as number | null, Validators.required],
        professorId: [null as number | null, Validators.required],
        requiredRoomTypeId: [null as number | null, Validators.required],
        weeklySessions: [1, [Validators.required, Validators.min(1), Validators.max(10)]],
        isSplitPerGroup: [false]
    });

    // Computed: Filter student groups by selected major chip
    visibleGroups = computed(() => {
        const majorId = this.selectedMajorFilter();
        const all = this.groups();
        if (!majorId) return all;
        return all.filter(g => g.levelId === majorId);
    });

    // Computed: Map blocks by subjectId for fast matrix rendering
    blocksBySubject = computed(() => {
        const map = new Map<number, CourseOfferingBlock[]>();
        for (const b of this.blocks()) {
            const existing = map.get(b.subjectId) ?? [];
            existing.push(b);
            map.set(b.subjectId, existing);
        }
        return map;
    });

    // Group saved blocks into unique Subject Columns (Same subject -> same column, no duplicates)
    subjectColumns = computed<SubjectColumn[]>(() => {
        const colMap = new Map<number, SubjectColumn>();
        const majorId = this.selectedMajorFilter();
        const visibleGroupSet = new Set(this.visibleGroups().map(g => g.id));

        for (const b of this.blocks()) {
            if (majorId && !b.groups.some(g => visibleGroupSet.has(g.groupId))) {
                continue;
            }

            let col = colMap.get(b.subjectId);
            if (!col) {
                col = {
                    subjectId: b.subjectId,
                    subjectCode: b.subjectCode ?? '',
                    subjectName: b.subjectName ?? '',
                    blocks: []
                };
                colMap.set(b.subjectId, col);
            }

            if (!col.blocks.some(existing => existing.sessionUuid === b.sessionUuid)) {
                col.blocks.push(b);
            }
        }

        return Array.from(colMap.values());
    });

    ngOnInit(): void {
        this.shellSrv.setTitle('COURSE_OFFERINGS');
        this.shellSrv.setGoBack(this.exit.bind(this));
        const id = Number(this.route.snapshot.paramMap.get('termId') || 0);
        this.termId.set(id);
        if (!id) {
            this.exit()
            return;
        }
        this.loadTermAndGroups(id);
        this.loadBlocks();
        this.loadReferenceData();
    }

    private loadTermAndGroups(termId: number): void {
        this.termSrv.get(termId).subscribe({
            next: (t) => {
                this.term.set(t);
                if (t?.academicYearId) {
                    this.groupSrv.getByYearId(t.academicYearId).subscribe({
                        next: (res) => this.groups.set(res.data ?? [])
                    });
                }
            },
            error: err => {
                throw err;
            }
        });
    }

    private loadReferenceData(): void {
        this.majorSrv.getAll('', 'id').subscribe(res => this.majors.set(res.data ?? []));
        this.roomTypeSrv.getAll('', 'id').subscribe(res => this.roomTypes.set(res.data ?? []));
    }

    loadBlocks(): void {
        this.loading.set(true);
        this.offeringSrv.getBlocks(this.termId()).subscribe({
            next: (data) => {
                this.blocks.set(data ?? []);
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    // --- Interactive Board Actions ---

    openCreateBuilder(preselectedSubjectId?: number): void {
        this.editingSessionUuid.set(null);
        this.selectedGroupIds.set([]);
        this.builderForm.reset({
            subjectId: preselectedSubjectId ?? null,
            professorId: null,
            requiredRoomTypeId: this.roomTypes()[0]?.id ?? null,
            weeklySessions: 1,
            isSplitPerGroup: false
        });
        this.isBuilderOpen.set(true);
    }

    // TODO block could be null
    openEditBuilder(block: CourseOfferingBlock): void {
        this.editingSessionUuid.set(block.sessionUuid);
        this.selectedGroupIds.set(block.groups.map(g => g.groupId));
        this.builderForm.patchValue({
            subjectId: block.subjectId,
            professorId: block.professorId,
            requiredRoomTypeId: block.requiredRoomTypeId,
            weeklySessions: block.weeklySessions,
            isSplitPerGroup: false
        });
        this.isBuilderOpen.set(true);
    }

    closeBuilder(): void {
        this.isBuilderOpen.set(false);
        this.editingSessionUuid.set(null);
        this.selectedGroupIds.set([]);
    }

    toggleGroupSelection(groupId: number): void {
        const current = this.selectedGroupIds();
        if (current.includes(groupId)) {
            this.selectedGroupIds.set(current.filter(id => id !== groupId));
        } else {
            this.selectedGroupIds.set([...current, groupId]);
        }
    }

    saveSessionBlock(): void {
        const groupIds = this.selectedGroupIds();
        if (this.builderForm.invalid || groupIds.length === 0) return;

        const formVal = this.builderForm.value;
        const payload: CourseOfferingBatchRequest = {
            termId: this.termId(),
            subjectId: formVal.subjectId!,
            professorId: formVal.professorId!,
            requiredRoomTypeId: formVal.requiredRoomTypeId!,
            studentGroupIds: groupIds,
            weeklySessions: formVal.weeklySessions ?? 1,
            isSplitPerGroup: Boolean(formVal.isSplitPerGroup)
        };

        const editUuid = this.editingSessionUuid();
        const request$ = editUuid
            ? this.offeringSrv.updateBatch(editUuid, payload)
            : this.offeringSrv.createBatch(payload);

        request$.subscribe({
            next: () => {
                this.closeBuilder();
                this.loadBlocks();
            }
        });
    }

    deleteBlock(sessionUuid: string, event: Event): void {
        event.stopPropagation();
        this.offeringSrv.deleteBatch(sessionUuid).subscribe({
            next: () => {
                if (this.editingSessionUuid() === sessionUuid) {
                    this.closeBuilder();
                }
                this.loadBlocks();
            }
        });
    }

    removeSingleGroupFromBlock(offeringId: number, event: Event): void {
        event.stopPropagation();
        this.offeringSrv.deleteSingleRow(offeringId).subscribe({
            next: () => this.loadBlocks()
        });
    }

    exit(): void {
        this.router.navigate(['..'], { relativeTo: this.route }).then();
    }

    onDockGroupClick(groupId: number): void {
        if (!this.isBuilderOpen()) {
            this.openCreateBuilder();
        }
        this.toggleGroupSelection(groupId);
    }

    goToTimetable(): void {
        this.router.navigate(['timetable'], { relativeTo: this.route }).then();
    }
}
