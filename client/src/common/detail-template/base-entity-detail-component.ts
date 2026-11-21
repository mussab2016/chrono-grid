import { Directive, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { BaseEntityService } from '../service-template/service-template.component';
import { ShellService } from '../shell-service/shell.service';

@Directive()
export abstract class BaseEntityDetailComponent<T extends { id: number }> implements OnInit, OnDestroy {
    // Concrete service provided by child components
    public abstract entitySrv: BaseEntityService<T>;

    // Shared core services via Angular inject API
    public shellSrv = inject(ShellService);
    public route = inject(ActivatedRoute);
    public router = inject(Router);
    public formBuilder = inject(FormBuilder);

    // Reactive signal holding the FormGroup; updates UI automatically when HTTP data arrives
    public form = signal<FormGroup | null>(null);
    protected entityId = 0;

    public ngOnInit(): void {
        this.shellSrv.setGoBack(this.exit.bind(this));

        const paramId = this.route.snapshot.paramMap.get('id');
        this.entityId = paramId ? +paramId : 0;
        this.form.set(null);

        if (this.entityId > 0) {
            this.entitySrv.get(this.entityId).subscribe({
                next: (entity) => {
                    this.form.set(this.setupForm(entity));
                },
                error: (err) => {
                    throw err;
                }
            });
        } else {
            this.form.set(this.setupForm(null));
        }
    }

    protected readForm(): T {
        return this.form()?.getRawValue() as T;
    }

    public save(): void {
        const currentForm = this.form();
        if (!currentForm || currentForm.invalid) {
            currentForm?.markAllAsTouched();
            return;
        }

        this.entitySrv.save(this.readForm()).subscribe({
            next: () => {
                this.exit();
            },
            error: (err) => {
                throw err;
            }
        });
    }

    public exit(): void {
        this.router.navigate(['..'], { relativeTo: this.route }).then();
    }

    public ngOnDestroy(): void {
        this.shellSrv.setTitle('');
        this.shellSrv.setGoBack(undefined);
    }

    // Must be implemented by child components to build the specific form controls
    public abstract setupForm(entity: T | null): FormGroup;
}
