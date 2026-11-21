import {inject, Injectable} from '@angular/core';
import {Direction} from '@angular/cdk/bidi';
import {HttpClient} from "@angular/common/http";
import {Observable, Subject} from "rxjs";

export interface Language {
    code: 'ar' | 'en' | 'fr';
    dir: Direction;
    name: string;
}

export const DEFAULT_SUPPORTED_LANGUAGES: Language[] = [
    {code: 'ar', dir: 'rtl', name: 'العربية'},
    {code: 'en', dir: 'ltr', name: 'English'},
    {code: 'fr', dir: 'ltr', name: 'Français'}
];

const DEFAULT_LANGUAGE = DEFAULT_SUPPORTED_LANGUAGES[2]; // en

@Injectable({
    providedIn: 'root'
})

export class LanguageService {

    private readonly LANG_TOKEN = 'language';

    public rtl = true;
    public dir: Direction = 'ltr';
    private currentLanguage: Language = DEFAULT_LANGUAGE;
    private dictionary: Record<string, string> = {};
    private changeSubject = new Subject<string>();
    private http = inject(HttpClient);

    constructor() {
        // Initialize default language immediately when service is created
        this.loadInitialLanguage();
    }

    public loadInitialLanguage() {
        const code = localStorage.getItem(this.LANG_TOKEN) as 'ar' | 'en' | 'fr' | null;
        const lang = DEFAULT_SUPPORTED_LANGUAGES.find(l => l.code === code) || DEFAULT_LANGUAGE;
        this.setLanguage(lang);
    }

    public get_dir(): Direction {
        return this.currentLanguage.dir;
    }

    // public getCurrentLanguage(): string {
    //     return this.currentLanguage?.code;
    // }

    setLanguage(language: Language): void {
        this.currentLanguage = language;
        localStorage.setItem(this.LANG_TOKEN, language.code);
        this.dir = language.dir;
        this.rtl = this.dir == 'rtl';
        this.loadLanguage(language.code).subscribe({
            next: (dict) => {
                this.dictionary = dict;
                this.changeSubject.next(language.code);
            },
            error: (err) => console.error('Error loading translation:', err)
        });
    }

    loadLanguage(lang: 'ar' | 'en' | 'fr') {
        const url = `translate/${lang}.json`;
        return this.http.get<Record<string, string>>(url);
    }
    //
    // private get_defaultLanguage() {
    //     const code: string = !!localStorage.getItem('language') ? <string>localStorage.getItem('language') : DEFAULT_LANGUAGE;
    //     this.currentLanguage = <Language>DEFAULT_SUPPORTED_LANGUAGES.find(l => l.code === code);
    //     //    this.changeSubject.next(this.currentLanguage.code)
    // }

    translate(key: string): string {
        return this.dictionary[key] || key;
    }

    public translateAll(keys: string[]) {
        return keys.map(k => this.translate(k));
    }

    public get languageChange(): Observable<string> {
        return this.changeSubject.asObservable();
    }

}
