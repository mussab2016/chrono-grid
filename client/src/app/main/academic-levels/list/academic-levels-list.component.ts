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
import {AcademicLevel} from '../../../api/academic-level';
import {AcademicLevelService} from '../academic-level.service';
import {MatSort, MatSortHeader} from '@angular/material/sort';
import {EntityCardSchema} from '@common/entity-cards/entity-card.model';
import {EntityCardComponent} from '@common/entity-cards/entity-card.component';


@Component({
    selector: 'academic-levels-list',
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
    templateUrl: './academic-levels-list.component.html',
    styleUrl: './academic-levels-list.component.scss'
})
export class AcademicLevelsListComponent extends BaseEntityListComponent<AcademicLevel> implements OnInit {

    public override entitySrv = inject(AcademicLevelService);

    public displayedColumns: string[] = ['code', 'name'];


    public cardSchema: EntityCardSchema<AcademicLevel> = {
        titleKey: 'code',
        subtitleKey: 'name',
        subtitleIcon: 'fa-solid fa-graduation-cap',
        fields: [
            { key: 'id', label: 'ID', icon: 'fa-regular fa-id-badge', isBadge: true }
        ]
    };
}
