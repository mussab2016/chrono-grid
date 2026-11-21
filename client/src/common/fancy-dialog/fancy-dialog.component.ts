// fancy-dialog.component.ts
import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import {TranslateDirective} from "../translate-directive/translate.directive";

export type FancyDialogType = 'confirm' | 'error';

export interface FancyDialogData {
    type: FancyDialogType;
    msg: string;
}

@Component({
    selector: 'fancy-dialog',
    standalone: true,
    imports: [CommonModule, MatDialogModule, MatButtonModule, TranslateDirective],
    templateUrl:'fancy-dialog.component.html',
    styleUrls: ['./fancy-dialog.component.scss']
})
export class FancyDialogComponent {
    constructor(
        @Inject(MAT_DIALOG_DATA) public data: FancyDialogData,
        private dialogRef: MatDialogRef<FancyDialogComponent>
    ) {}

    close(result?: boolean) {
        this.dialogRef.close(result);
    }
}
