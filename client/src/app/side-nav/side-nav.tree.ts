import {TimetableComponent} from '../main/timetable-export/model/timetable.component';

export interface SideNavNode {
    name: string;
    icon: string;
    route: string;
}

const MAIN_ROUTE = 'main'

export const PROFESSORS_ROUTE = 'professors';
export const ROOM_TYPES_ROUTE = 'room-types';
export const ROOMS_ROUTE = 'rooms';
export const ACADEMIC_LEVELS_ROUTE = 'academic-levels';
export const SUBJECTS_ROUTE = 'subjects';
export const ACADEMIC_CALENDAR_ROUTE = 'academic-calendar';
export const TIMETABLE_EXPORT_COMPONENT_ROUTE = 'timetable';
// export const COURSE_OFFERING_ROUTE = 'course-offering';

const joinRoute = (...segments: string[]): string => segments.join('/');

export const FULL_PROFESSORS_ROUTE = joinRoute(MAIN_ROUTE, PROFESSORS_ROUTE);
export const FULL_ROOM_TYPES_ROUTE = joinRoute(MAIN_ROUTE, ROOM_TYPES_ROUTE);
export const FULL_ROOMS_ROUTE = joinRoute(MAIN_ROUTE, ROOMS_ROUTE);
export const FULL_ACADEMIC_LEVELS_ROUTE = joinRoute(MAIN_ROUTE, ACADEMIC_LEVELS_ROUTE);
export const FULL_SUBJECTS_ROUTE = joinRoute(MAIN_ROUTE, SUBJECTS_ROUTE);
export const FULL_ACADEMIC_CALENDAR_ROUTE = joinRoute(MAIN_ROUTE, ACADEMIC_CALENDAR_ROUTE);
export const FULL_TIMETABLE_EXPORT_COMPONENT = joinRoute(MAIN_ROUTE, TIMETABLE_EXPORT_COMPONENT_ROUTE)
// export const FULL_COURSE_OFFERING_ROUTE = joinRoute(MAIN_ROUTE, COURSE_OFFERING_ROUTE);

export const ErpNavTree: SideNavNode[] = [
    {
        name: 'PROFESSORS',
        icon: 'fa-solid fa-chalkboard-user',//'fa-solid fa-user-graduate',
        route: FULL_PROFESSORS_ROUTE,
    },
    {
        name: 'ROOM_TYPES',
        icon: 'fa-solid fa-layer-group',
        route: FULL_ROOM_TYPES_ROUTE,
    },
    {
        name: 'ROOMS',
        icon: 'fa-solid fa-door-open',
        route: FULL_ROOMS_ROUTE,
    },
    {
        name: 'ACADEMIC_LEVELS',
        icon: 'fa-solid fa-graduation-cap',
        route: FULL_ACADEMIC_LEVELS_ROUTE,
    },
    {
        name: 'SUBJECTS',
        icon: 'fa-solid fa-book-open',
        route: FULL_SUBJECTS_ROUTE,
    },
    {
        name: 'ACADEMIC_TERMS',
        icon: 'fa-solid fa-calendars',
        route: FULL_ACADEMIC_CALENDAR_ROUTE,
    },
    {
        name: 'TIMETABLE',
        icon: 'fa-solid fa-calendar-clock',
        route: FULL_TIMETABLE_EXPORT_COMPONENT,
    },
];
