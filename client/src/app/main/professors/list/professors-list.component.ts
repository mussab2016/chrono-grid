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
import {ProfessorService} from '../professor.service';
import {MatSort, MatSortHeader} from '@angular/material/sort';
import {EntityCardSchema} from '@common/entity-cards/entity-card.model';
import {EntityCardComponent} from '@common/entity-cards/entity-card.component';

@Component({
    selector: 'app-professors-list',
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
    templateUrl: './professors-list.component.html',
    styleUrl: './professors-list.component.scss'
})
export class ProfessorsListComponent extends BaseEntityListComponent<Professor> implements OnInit {

    public override entitySrv = inject(ProfessorService);

    public displayedColumns: string[] = [
        'fullName',
        'department'
    ];

    override ngOnInit(): void {
        this.shellSrv.setTitle('PROFESSORS');
        super.ngOnInit();
    }

    public cardSchema: EntityCardSchema<Professor> = {
        titleKey: 'fullName',
        subtitleKey: 'department',
        subtitleIcon: 'fa-solid fa-building-columns',
        fields: [
            { key: 'id', label: 'ID', icon: 'fa-regular fa-id-badge', isBadge: true }
            // { key: 'weeklyHours', label: 'HOURS', icon: 'fa-regular fa-clock' }
        ]
    };
}
