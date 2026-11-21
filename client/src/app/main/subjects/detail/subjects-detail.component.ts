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
import { Subject } from '../../../api/subject';
import {SubjectService} from '../subject.service';


@Component({
    selector: 'app-subjects-detail',
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
    templateUrl: './subjects-detail.component.html',
    styleUrl: './subjects-detail.component.scss'
})
export class SubjectsDetailComponent extends BaseEntityDetailComponent<Subject> implements OnInit {
    // Inject concrete professor service
    public override entitySrv = inject(SubjectService);

    override ngOnInit(): void {
        this.shellSrv.setTitle('PROFESSOR');
        super.ngOnInit();
    }

    override setupForm(entity: Subject | null): FormGroup {
        return this.formBuilder.group({
            id: [entity?.id || 0],
            code: [entity?.code ?? '', [Validators.maxLength(50)]],
            name: [entity?.name ?? '', [Validators.required, Validators.maxLength(255)]]
        });
    }
}
