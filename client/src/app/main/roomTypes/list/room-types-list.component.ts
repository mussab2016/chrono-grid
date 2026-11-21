import {Component, inject, OnInit} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormField, MatInput, MatLabel } from '@angular/material/input';
import { MatIcon } from '@angular/material/icon';
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

import { BaseEntityListComponent } from '@common/list-template/base-entity-list-component';
import { PaginationComponent } from '@common/pagination/pagination.component';
import { TranslateDirective } from '@common/translate-directive/translate.directive';
import { TranslatePipe } from '../../pipes/unix-transtale.pipe';
import { Professor } from '../../../api/professor';
import {RoomTypeService} from '../room-type.service';
import {MatSort, MatSortHeader} from '@angular/material/sort';
import {EntityCardSchema} from '@common/entity-cards/entity-card.model';
import {EntityCardComponent} from '@common/entity-cards/entity-card.component';
import {RoomType} from '../../../api/room-type';

@Component({
    selector: 'app-room-types-list',
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
    templateUrl: './roomTypes-list.component.html',
    styleUrl: './roomTypes-list.component.scss'
})
export class RoomTypesListComponent extends BaseEntityListComponent<RoomType> implements OnInit {

    public override entitySrv = inject(RoomTypeService);

    public displayedColumns: string[] = [
        'code',
        'description'
    ];

    override ngOnInit(): void {
        this.shellSrv.setTitle('PROFESSORS');
        super.ngOnInit();
    }

    public cardSchema: EntityCardSchema<RoomType> = {
        titleKey: 'code',
        subtitleKey: 'description',
        subtitleIcon: 'fa-solid fa-door-open',
        fields: [
            { key: 'id', label: 'ID', icon: 'fa-regular fa-id-badge', isBadge: true }
        ]
    };
}
