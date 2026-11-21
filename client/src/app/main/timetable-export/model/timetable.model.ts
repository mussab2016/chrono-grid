import { TimetableSlot } from '../../../api/timetable';
import * as XLSX from 'xlsx-js-style';

export type TimetableViewType = 'AcademicLevels' | 'Class' | 'Professor';

export interface EntityOption {
    id: number;
    label: string;
}

// Each line inside a cell carries numeric IDs so filtering is purely numeric and instant
export interface CellLineItem {
    text: string;
    academicLevelId?: number;
    groups?: number[];
}

export interface UnifiedTableRow {
    id: number;            // Always numeric (dayNumber, roomId, or professorId)
    label: string;
    cells: CellLineItem[][];
}

export interface UnifiedTimetable {
    cornerLabel: string;
    columnLabels: string[];
    rows: UnifiedTableRow[];
}

export const DAYS: readonly string[] = [
    'Saturday',
    'Sunday',
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday'
];

// Up to 10 default session slots (starting at 08:00, 1.5h each, 0 break)
export const SESSIONS: readonly string[] = [
    '08:00-09:30',
    '09:30-11:00',
    '11:00-12:30',
    '12:30-14:00',
    '14:00-15:30',
    '15:30-17:00',
    '17:00-18:30',
    '18:30-20:00',
    '20:00-21:30',
    '21:30-23:00'
];

/**
 * Resolves the active session labels based on max slotNumber present in slots (capped at 10).
 */
export function resolveActiveSessions(
    slots: readonly TimetableSlot[],
    savedSessions: readonly string[] = SESSIONS
): string[] {
    let maxSlot = 0;
    for (const s of slots) {
        if (s.slotNumber > maxSlot) {
            maxSlot = s.slotNumber;
        }
    }
    const count = Math.min(10, maxSlot);
    if (count === 0) return [];

    return Array.from(
        { length: count },
        (_, idx) => savedSessions[idx] ?? SESSIONS[idx] ?? `Slot ${idx + 1}`
    );
}

/**
 * Resolves active day labels from DAYS using: startIndex = 1 - (daysCount - 5)
 * - 5 days -> starts at index 1 ('Sunday' .. 'Thursday')
 * - 6 days -> starts at index 0 ('Saturday' .. 'Thursday')
 */
export function resolveActiveDays(slots: readonly TimetableSlot[]): string[] {
    let maxDay = 0;
    for (const s of slots) {
        if (s.dayNumber > maxDay) {
            maxDay = s.dayNumber;
        }
    }

    const daysCount = maxDay >= 6 ? 6 : 5;
    const startIndex = Math.max(0, 1 - (daysCount - 5));
    return DAYS.slice(startIndex, startIndex + daysCount);
}

// ============================================================================
// A. FULL BUILDERS (Run once when raw slots or session labels change)
// ============================================================================

/**
 * 1. Builds the FULL Academic Levels timetable (all levels, all groups).
 * - Columns: Active sessions
 * - Rows: Active days
 */
export function buildFullAcademicLevelTimetable(
    slots: readonly TimetableSlot[],
    savedSessions: readonly string[] = SESSIONS
): UnifiedTimetable {
    const columnLabels = resolveActiveSessions(slots, savedSessions);
    const dayLabels = resolveActiveDays(slots);

    const grid: CellLineItem[][][] = Array.from({ length: dayLabels.length }, () =>
        Array.from({ length: columnLabels.length }, () => [])
    );

    for (const s of slots) {
        const dIdx = s.dayNumber - 1;
        const sIdx = s.slotNumber - 1;
        if (dIdx >= 0 && dIdx < dayLabels.length && sIdx >= 0 && sIdx < columnLabels.length) {
            const groupsText = s.groups.map(g => `G${g}`).join(', ');
            grid[dIdx][sIdx].push({
                academicLevelId: s.academicLevelId,
                groups: [...s.groups],
                text: `[${s.academicLevelCode}] ${s.subjectCode} (${s.roomTypeCode}) - ${s.professorName} | ${s.roomCode} (${groupsText})`
            });
        }
    }

    const rows: UnifiedTableRow[] = dayLabels.map((dayName, dIdx) => ({
        id: dIdx + 1,
        label: dayName,
        cells: grid[dIdx]
    }));

    return {
        cornerLabel: 'Day / Time',
        columnLabels,
        rows
    };
}

/**
 * 2. Builds the FULL Classrooms timetable (all server rooms in exact backend order).
 * - Columns: Active days
 * - Rows: Server rooms
 */
export function buildFullClassTimetable(
    slots: readonly TimetableSlot[],
    serverRooms: readonly EntityOption[],
    savedSessions: readonly string[] = SESSIONS
): UnifiedTimetable {
    const sessionLabels = resolveActiveSessions(slots, savedSessions);
    const dayLabels = resolveActiveDays(slots);

    const roomIndexMap = new Map<number, number>();
    serverRooms.forEach((r, idx) => roomIndexMap.set(r.id, idx));

    const grid: (CellLineItem | null)[][][] = Array.from({ length: serverRooms.length }, () =>
        Array.from({ length: dayLabels.length }, () =>
            Array.from({ length: sessionLabels.length }, () => null)
        )
    );

    for (const s of slots) {
        const rIdx = roomIndexMap.get(s.roomId);
        const dIdx = s.dayNumber - 1;
        const sIdx = s.slotNumber - 1;
        if (
            rIdx !== undefined &&
            dIdx >= 0 &&
            dIdx < dayLabels.length &&
            sIdx >= 0 &&
            sIdx < sessionLabels.length
        ) {
            const timeRange = sessionLabels[sIdx] ?? `Slot ${s.slotNumber}`;
            grid[rIdx][dIdx][sIdx] = {
                text: `${s.professorName} (${timeRange})`
            };
        }
    }

    const rows: UnifiedTableRow[] = serverRooms.map((room, rIdx) => ({
        id: room.id,
        label: room.label,
        cells: grid[rIdx].map(daySlots =>
            daySlots.filter((item): item is CellLineItem => item !== null)
        )
    }));

    return {
        cornerLabel: 'Classroom / Day',
        columnLabels: dayLabels,
        rows
    };
}

/**
 * 3. Builds the FULL Professors timetable (all server professors in exact backend order).
 * - Columns: Active days
 * - Rows: Server professors
 */
export function buildFullProfessorTimetable(
    slots: readonly TimetableSlot[],
    serverProfessors: readonly EntityOption[],
    savedSessions: readonly string[] = SESSIONS
): UnifiedTimetable {
    const sessionLabels = resolveActiveSessions(slots, savedSessions);
    const dayLabels = resolveActiveDays(slots);

    const profIndexMap = new Map<number, number>();
    serverProfessors.forEach((p, idx) => profIndexMap.set(p.id, idx));

    const grid: (CellLineItem | null)[][][] = Array.from({ length: serverProfessors.length }, () =>
        Array.from({ length: dayLabels.length }, () =>
            Array.from({ length: sessionLabels.length }, () => null)
        )
    );

    for (const s of slots) {
        const pIdx = profIndexMap.get(s.professorId);
        const dIdx = s.dayNumber - 1;
        const sIdx = s.slotNumber - 1;
        if (
            pIdx !== undefined &&
            dIdx >= 0 &&
            dIdx < dayLabels.length &&
            sIdx >= 0 &&
            sIdx < sessionLabels.length
        ) {
            const timeRange = sessionLabels[sIdx] ?? `Slot ${s.slotNumber}`;
            const groupsStr = s.groups.map(g => `G${g}`).join(', ');
            grid[pIdx][dIdx][sIdx] = {
                text: `${timeRange} | ${s.roomCode} | ${s.academicLevelCode} (${groupsStr})`
            };
        }
    }

    const rows: UnifiedTableRow[] = serverProfessors.map((prof, pIdx) => ({
        id: prof.id,
        label: prof.label,
        cells: grid[pIdx].map(daySlots =>
            daySlots.filter((item): item is CellLineItem => item !== null)
        )
    }));

    return {
        cornerLabel: 'Professor / Day',
        columnLabels: dayLabels,
        rows
    };
}

// ============================================================================
// B. FAST NUMERIC FILTERS (Operate directly on the cached UnifiedTimetable)
// ============================================================================

/**
 * Filters a cached Academic Level timetable by numeric levelId and optional numeric groupId.
 * If both are null/undefined, returns the cached full table as-is.
 */
export function filterAcademicLevelTimetable(
    fullTable: UnifiedTimetable,
    levelId: number | null,
    groupId: number | null
): UnifiedTimetable {
    if (levelId === null && groupId === null) {
        return fullTable;
    }

    return {
        ...fullTable,
        rows: fullTable.rows.map(row => ({
            ...row,
            cells: row.cells.map(cellItems =>
                cellItems.filter(item => {
                    const matchLevel = levelId === null || item.academicLevelId === levelId;
                    const matchGroup =
                        groupId === null || (item.groups !== undefined && item.groups.includes(groupId));
                    return matchLevel && matchGroup;
                })
            )
        }))
    };
}

/**
 * Filters a cached Classrooms timetable by numeric roomId.
 * If roomId is null, returns the cached full table as-is.
 */
export function filterClassTimetable(
    fullTable: UnifiedTimetable,
    roomId: number | null
): UnifiedTimetable {
    if (roomId === null) {
        return fullTable;
    }
    return {
        ...fullTable,
        rows: fullTable.rows.filter(row => row.id === roomId)
    };
}

/**
 * Filters a cached Professors timetable by numeric professorId.
 * If professorId is null, returns the cached full table as-is.
 */
export function filterProfessorTimetable(
    fullTable: UnifiedTimetable,
    professorId: number | null
): UnifiedTimetable {
    if (professorId === null) {
        return fullTable;
    }
    return {
        ...fullTable,
        rows: fullTable.rows.filter(row => row.id === professorId)
    };
}

/**
 * Exports any UnifiedTimetable to a styled Excel (.xlsx) file
 * dynamically matching the active Angular Material 3 theme tokens,
 * with generous cell padding, row heights, and column widths to prevent text overlap.
 */
export function exportTimetableToExcel(
    table: UnifiedTimetable,
    fileName: string = 'timetable.xlsx'
): void {
    // 1. Clean text without artificial space padding
    const headerRow: string[] = [table.cornerLabel, ...table.columnLabels];

    const bodyRows: string[][] = table.rows.map((row) => [
        row.label,
        ...row.cells.map((cellItems) =>
            cellItems.length > 0
                ? cellItems.map((item) => item.text).join('\n')
                : '—'
        )
    ]);

    const aoa: string[][] = [headerRow, ...bodyRows];
    const worksheet = XLSX.utils.aoa_to_sheet(aoa);

    // Extract live Material 3 theme colors dynamically
    const palette = {
        primary: resolveCssTokenToRgbHex('var(--mat-sys-primary)', '3F51B5'),
        onPrimary: resolveCssTokenToRgbHex('var(--mat-sys-on-primary)', 'FFFFFF'),
        secondaryContainer: resolveCssTokenToRgbHex('var(--mat-sys-secondary-container)', 'E8EAF6'),
        onSecondaryContainer: resolveCssTokenToRgbHex('var(--mat-sys-on-secondary-container)', '1E293B'),
        surfaceContainer: resolveCssTokenToRgbHex('var(--mat-sys-surface-container-high)', 'E2E8F0'),
        surface: resolveCssTokenToRgbHex('var(--mat-sys-surface)', 'FFFFFF'),
        onSurface: resolveCssTokenToRgbHex('var(--mat-sys-on-surface)', '1E293B'),
        onSurfaceVariant: resolveCssTokenToRgbHex('var(--mat-sys-on-surface-variant)', '64748B'),
        outlineVariant: resolveCssTokenToRgbHex('var(--mat-sys-outline-variant)', 'CBD5E1')
    };

    const borderStyle = {
        top: { style: 'thin', color: { rgb: palette.outlineVariant } },
        bottom: { style: 'thin', color: { rgb: palette.outlineVariant } },
        left: { style: 'thin', color: { rgb: palette.outlineVariant } },
        right: { style: 'thin', color: { rgb: palette.outlineVariant } }
    };

    const headerStyle = {
        font: { name: 'Calibri', sz: 13, bold: true, color: { rgb: palette.onPrimary } },
        fill: { fgColor: { rgb: palette.primary } },
        alignment: { horizontal: 'center', vertical: 'center', wrapText: true },
        border: borderStyle
    };

    const rowHeaderStyle = {
        font: { name: 'Calibri', sz: 12, bold: true, color: { rgb: palette.onSurface } },
        fill: { fgColor: { rgb: palette.surfaceContainer } },
        alignment: { horizontal: 'left', vertical: 'center', wrapText: true },
        border: borderStyle
    };

    const sessionCellStyle = {
        font: { name: 'Calibri', sz: 15, bold: true, color: { rgb: palette.onSecondaryContainer } },
        fill: { fgColor: { rgb: palette.secondaryContainer } },
        alignment: { horizontal: 'center', vertical: 'center', wrapText: true },
        border: borderStyle
    };

    const emptyCellStyle = {
        font: { name: 'Calibri', sz: 11, color: { rgb: palette.onSurfaceVariant } },
        fill: { fgColor: { rgb: palette.surface } },
        alignment: { horizontal: 'center', vertical: 'center' },
        border: borderStyle
    };

    const range = XLSX.utils.decode_range(worksheet['!ref'] ?? 'A1');
    for (let R = range.s.r; R <= range.e.r; ++R) {
        for (let C = range.s.c; C <= range.e.c; ++C) {
            const cellAddress = XLSX.utils.encode_cell({ r: R, c: C });
            const cell = worksheet[cellAddress];
            if (!cell) continue;

            if (R === 0) {
                cell.s = headerStyle;
            } else if (C === 0) {
                cell.s = rowHeaderStyle;
            } else if (cell.v && cell.v !== '—') {
                cell.s = sessionCellStyle;
            } else {
                cell.s = emptyCellStyle;
            }
        }
    }

    // 2. Exact proportional column widths (fit to longest line + tiny 2-char margin, capped at 32)
    const FONT_SCALE_FACTOR = 1.45;
    const MAX_COL_WCH = 48;
    const colWidths: number[] = headerRow.map((_, colIdx) => {
        let maxLen = colIdx === 0 ? 16 : 18;
        for (const r of aoa) {
            const cellVal = r[colIdx] ?? '';
            for (const line of cellVal.split('\n')) {
                const needed = line.length + 2;
                if (needed > maxLen) {
                    maxLen = Math.min(needed, MAX_COL_WCH);
                }
            }
        }
        return maxLen;
    });

    worksheet['!cols'] = colWidths.map((wch) => ({ wch }));

    // 3. Exact proportional row heights (16pt per actual line + 6pt margin)
    const CM_TO_PT = 28.3465;
    const EMPTY_ROW_HPT = Number((1.5 * CM_TO_PT).toFixed(2)); // 42.52 pt (1.5 cm)
    const VERTICAL_MARGIN_HPT = 1.0 * CM_TO_PT;               // 28.35 pt (0.5 cm top + 0.5 cm bottom)
    const LINE_HEIGHT_15PT = 20;                              // Actual height of 1 line at 15pt Bold
    worksheet['!rows'] = aoa.map((row, rowIdx) => {
        // Header row (13pt Bold + 1cm vertical margin)
        if (rowIdx === 0) {
            return { hpt: Math.round(18 + VERTICAL_MARGIN_HPT) };
        }

        // Check if all session cells in this row are empty
        const dataCells = row.slice(1);
        const isRowEmpty = dataCells.every((cellVal) => !cellVal || cellVal === '—');
        if (isRowEmpty) {
            return { hpt: EMPTY_ROW_HPT };
        }

        // Calculate max wrapped lines across ALL cells in the row (including col 0)
        let maxLines = 1;
        row.forEach((cellVal, colIdx) => {
            if (!cellVal || cellVal === '—') return;

            // Convert column's wch capacity back to actual 15pt characters that fit on 1 line
            const scale = colIdx === 0 ? 1.25 : FONT_SCALE_FACTOR;
            const charsPerLine = Math.max(8, Math.floor((colWidths[colIdx] - 2) / scale));

            let linesCount = 0;
            for (const rawLine of cellVal.split('\n')) {
                linesCount += Math.max(1, Math.ceil(rawLine.length / charsPerLine));
            }

            if (linesCount > maxLines) {
                maxLines = linesCount;
            }
        });

        return {
            hpt: Number((maxLines * LINE_HEIGHT_15PT + VERTICAL_MARGIN_HPT).toFixed(2))
        };
    });

    // 4. Fit all columns onto 1 landscape page wide when printing/viewing
    worksheet['!pageSetup'] = {
        orientation: 'landscape',
        fitToWidth: 1,
        fitToHeight: 0
    };

    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, 'Timetable');
    XLSX.writeFile(workbook, fileName.endsWith('.xlsx') ? fileName : `${fileName}.xlsx`);
}

/**
 * Resolves any CSS color / Material 3 token (including light-dark(), rgb(), hex)
 * at runtime via the browser's computed style and returns a 6-char uppercase RGB hex string.
 */
function resolveCssTokenToRgbHex(cssVarExpression: string, fallbackHex: string): string {
    if (typeof document === 'undefined') return fallbackHex;

    const probe = document.createElement('div');
    probe.style.color = cssVarExpression;
    probe.style.display = 'none';
    document.body.appendChild(probe);

    const computedColor = getComputedStyle(probe).color;
    document.body.removeChild(probe);

    const match = computedColor.match(/rgba?\(\s*(\d+),\s*(\d+),\s*(\d+)/i);
    if (!match) return fallbackHex;

    return [match[1], match[2], match[3]]
        .map((val) => Number(val).toString(16).padStart(2, '0'))
        .join('')
        .toUpperCase();
}
