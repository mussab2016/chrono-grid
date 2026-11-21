import { Component, inject, OnInit } from '@angular/core';
import { FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from '@angular/material/card';
import { MatToolbarModule } from '@angular/material/toolbar';

import { BaseEntityDetailComponent } from '@common/detail-template/base-entity-detail-component';
import { TranslateDirective } from '@common/translate-directive/translate.directive';
import { TranslatePipe } from '../../pipes/unix-transtale.pipe';
import { Professor } from '../../../api/professor';
import { ProfessorService } from '../professor.service';

@Component({
    selector: 'professors-detail-component',
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
    templateUrl: './professors-detail.component.html',
    styleUrl: './professors-detail.component.scss'
})
export class ProfessorsDetailComponent extends BaseEntityDetailComponent<Professor> implements OnInit {
    // Inject concrete professor service
    public override entitySrv = inject(ProfessorService);

    override ngOnInit(): void {
        this.shellSrv.setTitle('PROFESSOR');
        super.ngOnInit();
    }

    override setupForm(entity: Professor | null): FormGroup {
        return this.formBuilder.group({
            id: [entity?.id || 0],
            fullName: [entity?.fullName || '', Validators.required],
            department: [entity?.department || '', Validators.required],
        });
    }
}
