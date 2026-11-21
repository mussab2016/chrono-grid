import {Component, inject, OnDestroy, OnInit} from '@angular/core';

import {MatIconModule} from '@angular/material/icon';
import {MatTooltipModule} from '@angular/material/tooltip';
import {ShellService} from './shell.service';
import {ShellButton} from './shell-button.model';
import {MatButtonModule} from "@angular/material/button";
import {ThemeSelectorComponent} from "../../app/themes/theme-selector.component";
import {LanguageComponent} from "../../app/language/language.component";
import {TranslateDirective} from "../translate-directive/translate.directive";
import {Subscription} from "rxjs";
import {FancyDialogService} from "../fancy-dialog/fancy-dialog.service";
import {Router} from "@angular/router";
import {ThemeService} from "../../app/services/themeService";
import {TranslatePipe} from "../../app/main/pipes/unix-transtale.pipe";
import {LanguageService} from "../../app/services/languageService";

@Component({
    selector: 'app-shell',
    templateUrl: './shell.component.html',
    styleUrls: ['./shell.component.scss'],
    standalone: true,
    imports: [MatButtonModule, MatTooltipModule, MatIconModule, ThemeSelectorComponent, LanguageComponent, TranslateDirective, TranslatePipe]
})
export class ShellComponent implements OnInit, OnDestroy {
    buttons: ShellButton[] = [];
    public goBack: () => void = () => {
    };
    private subs: Subscription = new Subscription();
    private dialog = inject(FancyDialogService)
    private router = inject(Router);
    private themeSrv = inject(ThemeService);
    private langSrv = inject(LanguageService);

    constructor(public shellService: ShellService) {
    }

    ngOnInit(): void {
        // Subscribe to buttons changes
        this.subs.add(
            this.shellService.buttons$.subscribe(btns => {
                this.buttons = btns;
            })
        );

        // Update goBack reactively
        this.subs.add(
            this.shellService.goBack$.subscribe(fn => {
                this.goBack = fn || (() => {
                });
            })
        );
    }

    executeAction(btn: ShellButton) {
        btn.action?.();
    }

    ngOnDestroy(): void {
        this.subs.unsubscribe(); // cleanup subscriptions
    }

    isMainPage(): boolean {
        return this.router.url == '/main'
    }

    setDarkMode(){
        this.themeSrv.setDarkMode()
    }

    public rtl(){
        return this.langSrv.rtl
    }

}
