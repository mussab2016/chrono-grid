import {Injectable} from '@angular/core';
import {BaseEntityService, HttpListResponse} from "@common/service-template/service-template.component";
import {StudentGroup} from '../../api/student-group';
import {Observable} from 'rxjs';
import {HttpParams} from '@angular/common/http';
import {AcademicTerm} from '../../api/academic-term';

@Injectable({
    providedIn: 'root'
})
export class StudentGroupService extends BaseEntityService<StudentGroup> {
    protected constructor() {
        super('/student-groups');
    }

    getByYearId(yearId: number): Observable<HttpListResponse<StudentGroup>> {
        const params = new HttpParams().set('yearId', yearId.toString());
        return this.httpSrv.getList<StudentGroup>('/student-groups/', params);
    }

    getByMajorYearId(majorId: number, yearId: number): Observable<HttpListResponse<StudentGroup>> {
        const params = new HttpParams().set('yearId', yearId.toString());
        params.set('majorId', majorId);
        return this.httpSrv.getList<StudentGroup>('/student-groups', params);
    }

}
