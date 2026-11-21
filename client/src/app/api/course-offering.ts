export interface CourseOffering {
    id: number;
    sessionUuid: string;
    termId: number;
    subjectId: number;
    professorId: number;
    requiredRoomTypeId: number;
    studentGroupId: number;
    weeklySessions: number;
    createdAt?: string;

    termLabel?: string;
    subjectName?: string;
    subjectCode?: string;
    professorName?: string;
    roomTypeName?: string;
    studentGroupCode?: string;
    levelCode?: string;
}

export interface GroupSummary {
    offeringId: number;
    groupId: number;
    code?: string;
    levelCode?: string;
}

export interface CourseOfferingBlock {
    sessionUuid: string;
    termId: number;
    termLabel?: string;
    subjectId: number;
    subjectName?: string;
    subjectCode?: string;
    professorId: number;
    professorName?: string;
    requiredRoomTypeId: number;
    roomTypeName?: string;
    weeklySessions: number;
    groups: GroupSummary[];
}

export interface CourseOfferingBatchRequest {
    sessionUuid?: string;
    termId: number;
    subjectId: number;
    professorId: number;
    requiredRoomTypeId: number;
    studentGroupIds: number[];
    weeklySessions: number;
    isSplitPerGroup: boolean;
}

export const COURSE_OFFERING_URL: string = '/course-offerings/';
export const COURSE_OFFERING_BATCH_URL: string = '/course-offerings-batch/';
