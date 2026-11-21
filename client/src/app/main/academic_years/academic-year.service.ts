import { Injectable } from '@angular/core';
import { BaseEntityService } from '@common/service-template/service-template.component';
import { AcademicYear } from '../../api/academic-year';

@Injectable({ providedIn: 'root' })
export class AcademicYearService extends BaseEntityService<AcademicYear> {
    protected constructor() {
        super('/academic-years');
    }
}
