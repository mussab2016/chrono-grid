import {Injectable} from '@angular/core';
import {BaseEntityService} from "@common/service-template/service-template.component";
import {Subject} from '../../api/subject';

@Injectable({
    providedIn: 'root'
})
export class SubjectService extends BaseEntityService<Subject> {
    protected constructor() {
        super('/subjects');
    }
}
