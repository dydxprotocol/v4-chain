export namespace Commands {
    export { _commands_generate as generate };
    export { _commands_install as install };
    export { _commands_transpile as transpile };
}
export const Contracts: typeof Contracts;
import _commands_generate from "./commands/generate";
import _commands_install from "./commands/install";
import _commands_transpile from "./commands/transpile";
