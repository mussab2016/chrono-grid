import {Component, inject, OnInit} from '@angular/core';
import {FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {MatCardModule} from '@angular/material/card';
import {MatToolbarModule} from '@angular/material/toolbar';
import {BaseEntityDetailComponent} from '@common/detail-template/base-entity-detail-component';
import {TranslateDirective} from '@common/translate-directive/translate.directive';
import {TranslatePipe} from '../../pipes/unix-transtale.pipe';
import {Room} from '../../../api/room';
import {RoomsService} from '../rooms.service';
import {DefaultAutoCompleteModel} from '@common/ui-auto-complete/auto-complete-model';
import {RoomType} from '../../../api/room-type';
import {RoomTypeService} from '../../roomTypes/room-type.service';
import {UiAutoCompleteComponent} from '@common/ui-auto-complete/auto-complete.component';


@Component({
    selector: 'app-rooms-detail-component',
    standalone: true,
    imports: [
        ReactiveFormsModule,
        MatButtonModule,
        MatCardModule,
        MatToolbarModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        TranslateDirective,
        TranslatePipe,
        UiAutoCompleteComponent,
    ],
    templateUrl: './rooms-detail.component.html',
    styleUrl: './rooms-detail.component.scss'
})
export class RoomsDetailComponent extends BaseEntityDetailComponent<Room> implements OnInit {
    // Inject concrete professor service
    public override entitySrv = inject(RoomsService);
    private roomTypeSrv = inject(RoomTypeService);
    override ngOnInit(): void {
        this.shellSrv.setTitle('ROOM');
        super.ngOnInit();
    }

    override setupForm(entity: Room | null): FormGroup {
        return this.formBuilder.group({
            id: [entity?.id || 0],
            code: [entity?.code || '', [Validators.required, Validators.maxLength(100)]],
            type: [entity?.type || null, Validators.required],
            description: [entity?.description ?? null]
        });
    }

    roomTypeAutoComplete = new DefaultAutoCompleteModel<RoomType>(
        this.roomTypeSrv,
        'code',
        (rt) => rt.description ? `${rt.code} - ${rt.description}` : rt.code,
        'id'
    );
}
