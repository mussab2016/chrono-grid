// app/global-error.handler.ts
import { ErrorHandler, Injectable, inject } from '@angular/core';
import {FancyDialogService} from "../fancy-dialog/fancy-dialog.service";

@Injectable()
export class GlobalErrorHandler implements ErrorHandler {

    private dialog = inject(FancyDialogService); // inject the dialog service

    handleError(error: any): void {
        console.error('Global Error:', error);

        // Show confirmation or alert
        this.dialog.error(error); // You can translate or customize
    }
}
