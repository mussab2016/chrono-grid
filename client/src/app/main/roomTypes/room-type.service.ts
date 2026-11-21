import {Injectable} from '@angular/core';
import {BaseEntityService} from "@common/service-template/service-template.component";
import {RoomType} from '../../api/room-type';

@Injectable({
    providedIn: 'root'
})
export class RoomTypeService extends BaseEntityService<RoomType> {
    protected constructor() {
        super('/room-types');
    }
}
