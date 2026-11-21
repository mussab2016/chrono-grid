// fancy-dialog.service.ts
import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { FancyDialogComponent, FancyDialogData } from './fancy-dialog.component';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class FancyDialogService {
    constructor(private dialog: MatDialog) {}

    confirm(msg: string): Observable<boolean> {
        const dialogRef = this.dialog.open(FancyDialogComponent, {
            width: '400px',
            data: { type: 'confirm', msg } as FancyDialogData,
            panelClass: 'fancy-dialog-panel'
        });
        return dialogRef.afterClosed();
    }

    error(msg: string): Observable<void> {
        const dialogRef = this.dialog.open(FancyDialogComponent, {
            width: '400px',
            data: { type: 'error', msg } as FancyDialogData,
            panelClass: 'fancy-dialog-panel'
        });
        return dialogRef.afterClosed();
    }
}
