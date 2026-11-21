import {Injectable} from '@angular/core';
import {BaseEntityService} from "@common/service-template/service-template.component";
import {AcademicLevel} from '../../api/academic-level';

@Injectable({
    providedIn: 'root'
})
export class AcademicLevelService extends BaseEntityService<AcademicLevel> {
    protected constructor() {
        super('/academic-levels');
    }
}
