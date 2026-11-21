import {Component, inject, OnInit} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ReactiveFormsModule} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatFormField, MatInput, MatLabel} from '@angular/material/input';
import {MatIcon} from '@angular/material/icon';
import {
    MatCell,
    MatCellDef,
    MatColumnDef,
    MatHeaderCell,
    MatHeaderCellDef,
    MatHeaderRow,
    MatHeaderRowDef,
    MatRow,
    MatRowDef,
    MatTable
} from '@angular/material/table';
import {BaseEntityListComponent} from '@common/list-template/base-entity-list-component';
import {PaginationComponent} from '@common/pagination/pagination.component';
import {TranslateDirective} from '@common/translate-directive/translate.directive';
import {TranslatePipe} from '../../pipes/unix-transtale.pipe';
import {Room} from '../../../api/room';
import {RoomsService} from '../rooms.service';
import {MatSort, MatSortHeader} from '@angular/material/sort';
import {EntityCardSchema} from '@common/entity-cards/entity-card.model';
import {EntityCardComponent} from '@common/entity-cards/entity-card.component';

@Component({
    selector: 'app-rooms-list',
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatTable,
        MatColumnDef,
        MatHeaderCell,
        MatHeaderCellDef,
        MatCell,
        MatCellDef,
        MatHeaderRow,
        MatHeaderRowDef,
        MatRow,
        MatRowDef,
        PaginationComponent,
        TranslateDirective,
        TranslatePipe,
        MatFormField,
        MatLabel,
        MatIcon,
        MatInput,
        MatSort,
        MatSortHeader,
        EntityCardComponent
    ],
    templateUrl: './rooms-list.component.html',
    styleUrl: './rooms-list.component.scss'
})
export class RoomsListComponent extends BaseEntityListComponent<Room> implements OnInit {

    public override entitySrv = inject(RoomsService);

    public displayedColumns: string[] = [
        'code',
        'typeCode',
        'description'
    ];

    override ngOnInit(): void {
        this.shellSrv.setTitle('ROOMS');
        super.ngOnInit();
    }

    public cardSchema: EntityCardSchema<Room> = {
        titleKey: 'code',
        subtitleKey: 'description',
        subtitleIcon: 'fa-solid fa-door-open',
        fields: [
            {key: 'id', label: 'ID', icon: 'fa-regular fa-id-badge', isBadge: true},
            {key: 'type', label: 'ROOM_TYPE', icon: 'fa-solid fa-layer-group', isBadge: true}
        ]
    };
}
