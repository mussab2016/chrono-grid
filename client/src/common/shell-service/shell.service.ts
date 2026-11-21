import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
import {ShellButton} from './shell-button.model';

@Injectable({
    providedIn: 'root'
})
export class ShellService {
    private buttonsSubject = new BehaviorSubject<ShellButton[]>([]);
    buttons$ = this.buttonsSubject.asObservable();
    private _goBack: (() => void) | undefined = undefined;
    private _goBack$ = new BehaviorSubject<(() => void) | undefined>(undefined);
    private _title: string = '';

    setButtons(buttons: ShellButton[]) {
        this.buttonsSubject.next(buttons);
    }

    addButton(button: ShellButton) {
        const current = this.buttonsSubject.getValue();
        this.buttonsSubject.next([...current, button]);
    }

    clearButtons() {
        this.buttonsSubject.next([]);
    }

    goBack() {
        if (this._goBack) {
            return this._goBack;
        }
        return undefined
    }

    setGoBack(callback: (() => void) | undefined) {
        this._goBack = callback;
        this._goBack$.next(callback);
    }

    get goBack$() {
        return this._goBack$.asObservable();
    }

    setTitle(title: string) {
        this._title = title
    }

    getTitle() {
        return this._title
    }
}
