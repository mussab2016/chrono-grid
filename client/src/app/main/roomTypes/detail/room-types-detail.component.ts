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
import {RoomType} from '../../../api/room-type';
import {RoomTypeService} from '../room-type.service';

@Component({
    selector: 'Room-types-detail-component',
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
        TranslatePipe
    ],
    templateUrl: './room-types-detail.component.html',
    styleUrl: './room-types-detail.component.scss'
})
export class RoomTypesDetailComponent extends BaseEntityDetailComponent<RoomType> implements OnInit {
    // Inject concrete professor service
    public override entitySrv = inject(RoomTypeService);

    override ngOnInit(): void {
        this.shellSrv.setTitle('ROOM_TYPES');
        super.ngOnInit();
    }

    override setupForm(entity: RoomType | null): FormGroup {
        return this.formBuilder.group({
            id: [entity?.id || 0],
            code: [entity?.code || '', Validators.required],
            description: [entity?.description],
        });
    }
}
