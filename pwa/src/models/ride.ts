import ILocation from "./location";

export default interface IRide {
    id: string;
    type: string;
    date: number;
    destination: ILocation;
    direction: string;
    createdAt: number;
}