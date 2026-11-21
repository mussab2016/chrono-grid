export interface CardFieldConfig<T> {
    key: keyof T;
    label?: string;
    icon?: string;
    isBadge?: boolean;
}

export interface EntityCardSchema<T> {
    titleKey: keyof T;
    subtitleKey?: keyof T;
    subtitleIcon?: string;
    fields?: CardFieldConfig<T>[];
}
