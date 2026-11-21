import {Injectable} from '@angular/core';
import {
    COURSE_OFFERING_BATCH_URL,
    COURSE_OFFERING_URL,
    CourseOffering,
    CourseOfferingBatchRequest,
    CourseOfferingBlock,
} from '../../api/course-offering';
import {BaseEntityService} from '@common/service-template/service-template.component';
import {Observable} from 'rxjs';
import {HttpParams} from '@angular/common/http';


@Injectable({
    providedIn: 'root'
})
export class CourseOfferingService extends BaseEntityService<CourseOffering> {
    protected constructor() {
        super('/course-offerings');
    }

    // Fetch aggregated physical session blocks for a term (or a specific sessionUuid)
    getBlocks(termId?: number, sessionUuid?: string): Observable<CourseOfferingBlock[]> {
        let params = new HttpParams();
        if (termId && termId > 0) {
            params = params.set('termId', termId.toString());
        }
        if (sessionUuid) {
            params = params.set('sessionUuid', sessionUuid);
        }

        return this.httpSrv.getRawList<CourseOfferingBlock>(COURSE_OFFERING_BATCH_URL, params);
    }

    // Fetch raw flat rows if needed
    getRawByTerm(termId: number): Observable<CourseOffering[]> {
        const params = new HttpParams().set('termId', termId.toString());
        return this.httpSrv.getRawList<CourseOffering>(COURSE_OFFERING_URL, params);
    }

    // Create combined or split sessions for selected groups
    createBatch(payload: CourseOfferingBatchRequest): Observable<{ message: string }> {
        return this.httpSrv.post<{ message: string }>(COURSE_OFFERING_BATCH_URL, payload);
    }

    // Update an existing physical session block by sessionUuid
    updateBatch(sessionUuid: string, payload: CourseOfferingBatchRequest): Observable<{ message: string }> {
        return this.httpSrv.put<{ message: string }>(`${COURSE_OFFERING_BATCH_URL}${sessionUuid}`, payload);
    }


    // Delete an entire physical session block
    deleteBatch(sessionUuid: string): Observable<{ message: string }> {
        return this.httpSrv.remove<{ message: string }>(`${COURSE_OFFERING_BATCH_URL}${sessionUuid}`);
    }

    // Delete a single group row from a combined session
    deleteSingleRow(offeringId: number): Observable<{ message: string }> {
        return this.httpSrv.remove<{ message: string }>(`${COURSE_OFFERING_URL}${offeringId}`);
    }
}

