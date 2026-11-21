import {Injectable} from '@angular/core';
import {BaseEntityService} from "@common/service-template/service-template.component";
import {Room} from '../../api/room';

@Injectable({
    providedIn: 'root'
})
export class RoomsService extends BaseEntityService<Room> {
    protected constructor() {
        super('/rooms');
    }
}
