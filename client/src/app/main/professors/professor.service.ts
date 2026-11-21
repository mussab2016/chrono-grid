import {Injectable} from '@angular/core';
import {Professor} from "../../api/professor";
import {BaseEntityService} from "@common/service-template/service-template.component";

@Injectable({
    providedIn: 'root'
})
export class ProfessorService extends BaseEntityService<Professor> {
    protected constructor() {
        super('/professors');
    }
}
