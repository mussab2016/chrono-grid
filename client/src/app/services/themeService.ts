import {Injectable} from '@angular/core';
import {BehaviorSubject} from "rxjs";

@Injectable({
    providedIn: 'root'
})
export class ThemeService {

    private readonly THEME_KEY = 'app-theme';
    private readonly DEFAULT_THEME = 'theme-1';

    // Use BehaviorSubject to store current theme with initial value
    currentTheme = new BehaviorSubject<string>(this.DEFAULT_THEME);

    constructor() {
        this.getThemeFirstTime();
        this.setTheme(this.currentTheme.value); // apply saved or default theme
    }

    getThemeFirstTime() {
        const savedTheme = localStorage.getItem(this.THEME_KEY);
        const theme = savedTheme ? savedTheme : this.DEFAULT_THEME;
        this.currentTheme.next(theme);  // update BehaviorSubject
    }


    setTheme(theme: string) {
        const htmlBody = document.getElementById('html-body');
        if (htmlBody) {
            // Clean up old classes first (only themes)
            htmlBody.classList.forEach(c => {
                if (c.startsWith('theme-') || c.startsWith('dark-')) {
                    htmlBody.classList.remove(c);
                }
            });
            htmlBody.classList.add(theme);
        }
        localStorage.setItem(this.THEME_KEY, theme);
        this.currentTheme.next(theme); // 🔑 update subscribers
    }

    getCurrentThemeName(): string {
        return localStorage.getItem(this.THEME_KEY) ?? this.DEFAULT_THEME;
    }

    setDarkMode() {
        const theme = this.getCurrentThemeName();
        if (theme.startsWith('dark-')) {
            this.setTheme(theme.slice(5)) //  theme.replace(/^dark-/, ''); // removes only if it starts with dark- The ^ ensures it only removes dark- if it’s at the start.
        } else {
            this.setTheme('dark-' + theme)
        }
    }
}
