import { Component, ElementRef, inject, Input, ViewChild } from '@angular/core';
import { NgStyle } from '@angular/common';
import { MatIconButton, MatMiniFabButton } from '@angular/material/button';
import { MatMenuTrigger, MatMenu, MatMenuItem } from '@angular/material/menu';
import { MatInput } from '@angular/material/input';
import { LanguageService } from '../../app/services/languageService';
import { ExtendedPaginationModel, PaginationModel } from './pagination-model';

@Component({
    selector: 'ss-pagination',
    templateUrl: './pagination.component.html',
    styleUrls: ['./pagination.component.scss'],
    standalone: true,
    imports: [NgStyle, MatIconButton, MatMiniFabButton, MatMenuTrigger, MatInput, MatMenu, MatMenuItem]
})
export class PaginationComponent {
    public langSrv = inject(LanguageService);

    @Input()
    set paginationModel(model: PaginationModel) {
        this.unbindModel();
        this.model = model as ExtendedPaginationModel;
        this.bindModel();
        this.updateLastPage();
        this.syncInput();
    }

    model: ExtendedPaginationModel = new ExtendedPaginationModel(1, 15, 0);
    lastPage = 1;
    readonly pageSizes = [5, 10, 15, 20, 30, 40, 50];

    @ViewChild('pageInput', { static: false }) pageInput!: ElementRef<HTMLInputElement>;

    private bindModel(): void {
        if (!this.model) return;

        this.model.onPageSizeChange = () => {
            this.model.mPage = 1;
            this.updateLastPage();
            this.emitModelChangedEvent();
        };

        this.model.onElementCountChange = () => {
            this.updateLastPage();
            if (this.model.page > this.lastPage && this.lastPage > 0) {
                this.model.mPage = this.lastPage;
            }
            this.syncInput();
        };

        this.model.onPageChange = () => {
            this.syncInput();
            this.emitModelChangedEvent();
        };
    }

    private unbindModel(): void {
        if (!this.model) return;
        this.model.onPageChange = () => {};
        this.model.onPageSizeChange = () => {};
        this.model.onElementCountChange = () => {};
    }

    private updateLastPage(): void {
        if (!this.model || this.model.pageSize <= 0) {
            this.lastPage = 1;
            return;
        }
        this.lastPage = Math.max(1, Math.ceil(this.model.elementCount / this.model.pageSize));
    }

    private syncInput(): void {
        if (this.pageInput?.nativeElement) {
            this.pageInput.nativeElement.value = String(this.model.page);
        }
    }

    // Directional Navigation Handlers
    onFirstPage(): void {
        if (this.model.page !== 1) {
            this.model.page = 1;
        }
    }

    onLastPage(): void {
        if (this.model.page !== this.lastPage) {
            this.model.page = this.lastPage;
        }
    }

    onNextPage(): void {
        if (this.model.page < this.lastPage) {
            this.model.page++;
        }
    }

    onPrevPage(): void {
        if (this.model.page > 1) {
            this.model.page--;
        }
    }

    // RTL Aware actions
    onLeftButtonClick(): void {
        !this.langSrv.rtl ? this.onNextPage() : this.onPrevPage();
    }

    onLeftDoubleButtonClick(): void {
        !this.langSrv.rtl ? this.onLastPage() : this.onFirstPage();
    }

    onRightButtonClick(): void {
        !this.langSrv.rtl ? this.onPrevPage() : this.onNextPage();
    }

    onRightDoubleButtonClick(): void {
        !this.langSrv.rtl ? this.onFirstPage() : this.onLastPage();
    }

    isLeftDisabled(): boolean {
        return !this.langSrv.rtl ? this.model.page >= this.lastPage : this.model.page <= 1;
    }

    isRightDisabled(): boolean {
        return !this.langSrv.rtl ? this.model.page <= 1 : this.model.page >= this.lastPage;
    }

    onGoToPage(): void {
        if (!this.pageInput?.nativeElement) return;
        let val = parseInt(this.pageInput.nativeElement.value, 10);

        if (isNaN(val) || val < 1) {
            val = 1;
        } else if (val > this.lastPage) {
            val = this.lastPage;
        }

        this.model.page = val;
    }

    emitModelChangedEvent(): void {
        this.syncInput();
        if (this.model.onChange != null) {
            this.model.fireChangeEvent();
        }
    }

    // Filter numeric key events
    onInputKeyDown(evt: KeyboardEvent): boolean {
        const allowedKeys = ['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab', 'Enter'];
        if (allowedKeys.includes(evt.key) || (evt.ctrlKey && ['a', 'c', 'v', 'x'].includes(evt.key.toLowerCase()))) {
            return true;
        }
        return /^[0-9]$/.test(evt.key);
    }
}
