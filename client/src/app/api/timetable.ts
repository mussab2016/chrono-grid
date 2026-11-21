export const TIMETABLE_URL = '/timetables';
export const TIMETABLE_GENERATE_URL = '/timetables/generate';

export interface ProfPreferenceInput {
    professorId: number;
    dayNumber: number;
    slotNumber: number; // 0 = entire day
    score: number;      // -1 = forbidden, 10 = preferred
}

export interface GenerateTimetableRequest {
    termId: number;
    daysCount: number;
    slotsPerDay: number;
    preferences: ProfPreferenceInput[];
}

export interface TimetableSlot {
    id: number;
    termId: number;
    offeringId: number;
    sessionUuid: string;
    subjectId: number;
    subjectCode: string;
    subjectName: string;
    professorId: number;
    professorName: string;
    roomId: number;
    roomCode: string;
    roomTypeCode: string;
    academicLevelCode: string;
    academicLevelId: number;
    dayNumber: number;
    slotNumber: number;
    groups: number[];
}

export interface GenerateTimetableResponse {
    status: 'OPTIMAL' | 'FEASIBLE' | 'INFEASIBLE' | 'EMPTY_INPUT';
    slots: TimetableSlot[];
}
