import {Component, EventEmitter, inject, Output} from '@angular/core';
import {ThemeService} from "../services/themeService";
import {LanguageService} from "../services/languageService";

import {MatButtonModule} from '@angular/material/button';
import {MatMenuModule} from "@angular/material/menu";
import {TranslateDirective} from "@common/translate-directive/translate.directive";


@Component({
    selector: 'theme-selector',
    styleUrls: ['./theme-selector.component.scss'],
    templateUrl: './theme-selector.component.html',
    standalone: true,
    imports: [MatMenuModule, MatButtonModule, TranslateDirective]
})
export class ThemeSelectorComponent {
    private themeService = inject(ThemeService);
    public languageService = inject(LanguageService);

    @Output() onTheme = new EventEmitter<string>();

    themes: string[] = ['theme-1', 'theme-2', 'theme-3', 'theme-4', 'theme-5', 'theme-6', 'theme-7'];

    public setTheme(t: string) {
        this.themeService.setTheme(t);
    }
}
