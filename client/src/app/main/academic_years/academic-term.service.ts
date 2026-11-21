import {Injectable} from '@angular/core';
import {BaseEntityService, HttpListResponse} from "@common/service-template/service-template.component";
import {AcademicTerm} from '../../api/academic-term';
import {Observable} from 'rxjs';
import {HttpParams} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class AcademicTermService extends BaseEntityService<AcademicTerm> {
    protected constructor() {
        super('/academic-terms');
    }

    getByYearId(yearId: number): Observable<HttpListResponse<AcademicTerm>> {
        const params = new HttpParams().set('yearId', yearId.toString());
        return this.httpSrv.getList<AcademicTerm>('/academic-terms/', params);
    }
}
