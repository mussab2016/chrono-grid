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
import { AcademicLevelService } from '../academic-level.service';
import {AcademicLevel} from '../../../api/academic-level';

@Component({
    selector: 'academic-levels-detail-component',
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
    templateUrl: './academic-levels-detail.component.html',
    styleUrl: './academic-levels-detail.component.scss'
})
export class AcademicLevelsDetailComponent extends BaseEntityDetailComponent<AcademicLevel> implements OnInit {
    // Inject concrete professor service
    public override entitySrv = inject(AcademicLevelService);

    override ngOnInit(): void {
        this.shellSrv.setTitle('ACADEMIC_LEVEL');
        super.ngOnInit();
    }

    override setupForm(entity: AcademicLevel | null): FormGroup {
        return this.formBuilder.group({
            id: [entity?.id || 0],
            code: [entity?.code ?? '', [Validators.maxLength(10)]],
            name: [entity?.name ?? '', [Validators.required, Validators.maxLength(100)]]
        });
    }
}
