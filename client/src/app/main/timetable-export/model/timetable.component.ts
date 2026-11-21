import { Component, OnDestroy, computed, input, signal } from '@angular/core';
import { TimetableSlot } from '../../../api/timetable';
import {
    EntityOption,
    SESSIONS,
    TimetableViewType,
    UnifiedTimetable,
    buildFullAcademicLevelTimetable,
    buildFullClassTimetable,
    buildFullProfessorTimetable,
    filterAcademicLevelTimetable,
    filterClassTimetable,
    filterProfessorTimetable, exportTimetableToExcel
} from './timetable.model';

const SESSIONS_STORAGE_KEY = 'timetable_sessions_v1';
const MAX_SESSIONS_COUNT = 10;

@Component({
    selector: 'app-timetable',
    standalone: true,
    templateUrl: './timetable.component.html',
    styleUrl: './timetable.component.scss'
})
export class TimetableComponent implements OnDestroy {
    // Inputs from parent component
    readonly slots = input.required<readonly TimetableSlot[]>();
    readonly rooms = input<readonly EntityOption[]>([]);
    readonly professors = input<readonly EntityOption[]>([]);

    readonly type = input<TimetableViewType>('AcademicLevels');
    readonly targetId = input<number | null>(null);
    readonly groupFilter = input<number | null>(null);

    // Session labels persisted in localStorage (up to 10 slots)
    readonly sessionLabels = signal<string[]>(this.loadSavedSessions());
    readonly isSessionsDialogOpen = signal<boolean>(false);
    readonly draftSessions = signal<string[]>([]);

    // 1. Cached full tables (built once per slots/entities/sessions change)
    private readonly fullAcademicLevelTable = computed<UnifiedTimetable>(() =>
        buildFullAcademicLevelTimetable(this.slots(), this.sessionLabels())
    );

    private readonly fullClassTable = computed<UnifiedTimetable>(() =>
        buildFullClassTimetable(this.slots(), this.rooms(), this.sessionLabels())
    );

    private readonly fullProfessorTable = computed<UnifiedTimetable>(() =>
        buildFullProfessorTimetable(this.slots(), this.professors(), this.sessionLabels())
    );

    // 2. Active displayed table (fast numeric filtering over cached full tables)
    readonly displayedTable = computed<UnifiedTimetable>(() => {
        const viewType = this.type();
        const id = this.targetId();

        if (viewType === 'Class') {
            return filterClassTimetable(this.fullClassTable(), id);
        }

        if (viewType === 'Professor') {
            return filterProfessorTimetable(this.fullProfessorTable(), id);
        }

        return filterAcademicLevelTimetable(
            this.fullAcademicLevelTable(),
            id,
            this.groupFilter()
        );
    });

    // --- Sessions Configuration Dialog (up to 10 slots) ---
    openSessionsDialog(): void {
        const current = this.sessionLabels();
        const padded = Array.from(
            { length: MAX_SESSIONS_COUNT },
            (_, idx) => current[idx] ?? SESSIONS[idx] ?? ''
        );
        this.draftSessions.set(padded);
        this.isSessionsDialogOpen.set(true);
    }

    closeSessionsDialog(): void {
        this.isSessionsDialogOpen.set(false);
    }

    updateDraftSession(index: number, value: string): void {
        const next = [...this.draftSessions()];
        next[index] = value;
        this.draftSessions.set(next);
    }

    saveSessionsDialog(): void {
        const cleaned = this.draftSessions().map((val, idx) =>
            val.trim() ? val.trim() : (SESSIONS[idx] ?? `Slot ${idx + 1}`)
        );
        this.sessionLabels.set(cleaned);
        localStorage.setItem(SESSIONS_STORAGE_KEY, JSON.stringify(cleaned));
        this.isSessionsDialogOpen.set(false);
    }

    resetDefaultSessions(): void {
        this.draftSessions.set([...SESSIONS]);
    }

    ngOnDestroy(): void {
        // Component destroyed: computed signals and in-memory caches are automatically released
    }

    private loadSavedSessions(): string[] {
        try {
            const raw = localStorage.getItem(SESSIONS_STORAGE_KEY);
            if (raw) {
                const parsed = JSON.parse(raw);
                if (Array.isArray(parsed) && parsed.length > 0) {
                    return Array.from(
                        { length: MAX_SESSIONS_COUNT },
                        (_, idx) => parsed[idx] ?? SESSIONS[idx] ?? `Slot ${idx + 1}`
                    );
                }
            }
        } catch {
            // Ignore storage errors and use default SESSIONS
        }
        return [...SESSIONS];
    }

    exportExcel(): void {
        const table = this.displayedTable();
        const fileName = `Timetable_${this.type()}_${new Date().toISOString().slice(0, 10)}`;
        exportTimetableToExcel(table, fileName);
    }
}
