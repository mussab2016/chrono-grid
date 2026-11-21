export interface ShellButton {
    label: string;
    icon?: string;         // FontAwesome or MatIcon class
    matIcon?: string;      // If using Angular Material icon
    tooltip?: string;
    action: () => void;    // Function to execute on click
    color?: string;        // Optional button color
}
