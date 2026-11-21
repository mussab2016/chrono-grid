export interface Room {
    id: number;
    code: string;
    type: number;
    typeCode?: String; // The code of the room type
    description?: string;
    createdAt?: Date;
}
