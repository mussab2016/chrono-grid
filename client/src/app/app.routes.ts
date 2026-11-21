import {Routes} from '@angular/router';
import {MainComponent} from "./main/main.component";
// import {ModulesListComponent} from "./main/modules/modules.component";
import {ProfessorsListComponent} from "./main/professors/list/professors-list.component";
import {ProfessorsDetailComponent} from "./main/professors/detail/professors-detail.component";
import {RoomTypesListComponent} from './main/roomTypes/list/room-types-list.component';
import {RoomTypesDetailComponent} from './main/roomTypes/detail/room-types-detail.component';
import {
    ACADEMIC_CALENDAR_ROUTE,
    ACADEMIC_LEVELS_ROUTE,
    PROFESSORS_ROUTE,
    ROOM_TYPES_ROUTE,
    ROOMS_ROUTE,
    SUBJECTS_ROUTE, TIMETABLE_EXPORT_COMPONENT_ROUTE
} from './side-nav/side-nav.tree';
import {AcademicLevelsListComponent} from './main/academic-levels/list/academic-levels-list.component';
import {AcademicLevelsDetailComponent} from './main/academic-levels/detail/academic-levels-detail.component';
import {SubjectsListComponent} from './main/subjects/list/subjects-list.component';
import {SubjectsDetailComponent} from './main/subjects/detail/subjects-detail.component';
import {RoomsListComponent} from './main/rooms/list/rooms-list.component';
import {RoomsDetailComponent} from './main/rooms/detail/rooms-detail.component';
import {AcademicCalendarComponent} from './main/academic_years/component/academic-calendar.component';
import {CourseOfferingsComponent} from './main/course-offerings/component/course-offerings.component';
import {TimetableGeneratorComponent} from './main/timetable/component/timetable-generator.component';
import {TimetableComponent} from './main/timetable-export/model/timetable.component';
import {TimetablesExportComponent} from './main/timetable-export/timetables/timetables-export.component';

export const routes: Routes = [
    {
        path: 'main',
        canActivate: [],
        children: [
            {path: '', component: MainComponent},
            {
                path: PROFESSORS_ROUTE, children: [
                    {path: '', component: ProfessorsListComponent},
                    {path: ':id', component: ProfessorsDetailComponent}
                ]
            },
            {
                path: ROOM_TYPES_ROUTE, children: [
                    {path: '', component: RoomTypesListComponent},
                    {path: ':id', component: RoomTypesDetailComponent}
                ]
            },
            {
                path: ROOMS_ROUTE, children: [
                    {path: '', component: RoomsListComponent},
                    {path: ':id', component: RoomsDetailComponent}
                ]
            },
            {
                path: ACADEMIC_LEVELS_ROUTE, children: [
                    {path: '', component: AcademicLevelsListComponent},
                    {path: ':id', component: AcademicLevelsDetailComponent}
                ]
            },
            {
                path: SUBJECTS_ROUTE, children: [
                    {path: '', component: SubjectsListComponent},
                    {path: ':id', component: SubjectsDetailComponent}
                ]
            },
            {
                path: ACADEMIC_CALENDAR_ROUTE, children: [
                    {path: '', component: AcademicCalendarComponent},
                    {path: ':termId', component: CourseOfferingsComponent},
                    {path: ':termId/timetable', component: TimetableGeneratorComponent}
                ]
            },
            {
                path: TIMETABLE_EXPORT_COMPONENT_ROUTE, children: [
                    {path: '', component: TimetablesExportComponent},
                ]
            },
        ]
    },
    {path: '**', pathMatch: 'full', redirectTo: 'main'}
];
