import {Injectable} from '@angular/core';
import {BaseEntityService} from '@common/service-template/service-template.component';
import {
    GenerateTimetableRequest,
    GenerateTimetableResponse,
    TIMETABLE_GENERATE_URL,
    TIMETABLE_URL,
    TimetableSlot
} from '../../api/timetable';
import {Observable, timeout} from 'rxjs';
import {HttpParams} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class TimetableService extends BaseEntityService<TimetableSlot> {
    protected constructor() {
        super(TIMETABLE_URL);
    }

    // Invoke backend AI timetable optimization endpoint (60s timeout)
    generate(payload: GenerateTimetableRequest): Observable<GenerateTimetableResponse> {
        return this.httpSrv
            .post<GenerateTimetableResponse>(TIMETABLE_GENERATE_URL, payload)
            .pipe(timeout(600000));
    }

    // Retrieve persisted schedule slots for a given academic term
    getByTerm(termId: number): Observable<TimetableSlot[]> {
        const params = new HttpParams().set('termId', termId.toString());
        return this.httpSrv.getRawList<TimetableSlot>(TIMETABLE_URL, params);
    }
}
