import { Directive, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Sort } from '@angular/material/sort';
import { of, switchMap, debounceTime, distinctUntilChanged } from 'rxjs';

import { BaseEntityService } from '../service-template/service-template.component';
import { PaginatorPersistenceService } from '../pagination/paginator-persistence.service';
import { PaginationModel } from '../pagination/pagination-model';
import { FancyDialogService } from '../fancy-dialog/fancy-dialog.service';
import { ShellService } from '../shell-service/shell.service';

@Directive()
export abstract class BaseEntityListComponent<T extends { id: number }> implements OnInit, OnDestroy {
    // Services Injection (Modern Angular inject API)
    public abstract entitySrv: BaseEntityService<T>;
    public shellSrv = inject(ShellService);
    public route = inject(ActivatedRoute);
    public router = inject(Router);
    public pagPersistSrv = inject(PaginatorPersistenceService);
    private dialog = inject(FancyDialogService);

    // Pagination and Page Metadata
    public pagination: PaginationModel | undefined = PaginationModel.newInstance();
    public title: string = '';

    // Controls for Search and Sorting
    public filter = new FormControl<string>('', { nonNullable: true });

    // Signal State Management
    public data = signal<T[] | null>(null);
    public selected = signal<number>(0);
    public selectedRow = signal<T | null>(null);
    public viewMode = signal<'table' | 'cards'>('table');

    ngOnInit(): void {
        // Register shell back button action
        this.shellSrv.setGoBack(this.exit.bind(this));

        // Restore page persistence using current page title
        if (this.pagination) {
            this.pagination.copyModel(this.pagPersistSrv.get(this.title));
            if (!this.pagination.pageSize) {
                this.pagination.pageSize = 15;
            }
            this.pagination.onChange = () => this.refresh();
        }

        const savedSearch = this.pagPersistSrv.getSearch(this.title);
        this.filter.setValue(savedSearch, { emitEvent: false });

        // Setup reactive search debounce stream
        this.filter.valueChanges.pipe(
            debounceTime(350),
            distinctUntilChanged()
        ).subscribe(() => {
            if (this.pagination) {
                this.pagination.page = 1; // Reset to first page on search
            }
            this.refresh();
        });

        // Initial data fetch
        this.refresh();
    }

    public refresh(): void {
        this.data.set([]);
        this.select(0);

        of(null)
            .pipe(switchMap(() => this.getPage()))
            .subscribe({
                next: (resData) => {
                    this.data.set(resData);
                },
                error: (err) => {
                    this.data.set([]);
                    throw err;
                }
            });
    }

    public getPage() {
        return this.entitySrv.getPage(this.pagination!, this.filter.value);
    }

    public selectRow(row: T): void {
        this.selectedRow.set(row);
        const id = row ? row.id || 0 : 0;
        this.select(id);
    }

    public select(id: number): void {
        this.selected.set(id);
    }

    public edit(): void {
        if (this.selected()) {
            this.showDetail(this.selected());
        }
    }

    public add(): void {
        this.showDetail(0);
    }

    public showDetail(id: number): void {
        this.router.navigate([id], { relativeTo: this.route });
    }

    public exit(): void {
        this.router.navigate(['..'], { relativeTo: this.route });
    }

    public deleteSelected(): void {
        const currentSelected = this.selected();
        if (!currentSelected || currentSelected <= 0) return;

        this.dialog.confirm('DELETE_CONFIRMATION').subscribe((ok) => {
            if (ok) {
                const itemToDelete = this.data()?.find((x) => x.id === currentSelected);
                if (itemToDelete && itemToDelete.id) {
                    this.entitySrv.remove(itemToDelete.id).subscribe({
                        next: () => {
                            this.refresh();
                        },
                        error: (err) => {
                            throw err;
                        }
                    });
                }
            }
        });
    }

    public sort(event: Sort): void {
        this.entitySrv.sortColumn = event;
        this.refresh();
    }

    public toggleViewMode(mode: 'table' | 'cards'): void {
        this.viewMode.set(mode);
    }

    ngOnDestroy(): void {
        // Save state first before exiting
        if (this.pagination) {
            this.pagPersistSrv.set(this.title, this.pagination, this.filter.value);
        }
        this.shellSrv.setTitle('');
        this.shellSrv.setGoBack(undefined);
        this.entitySrv.sortColumn = { active: '', direction: '' };
    }
}
