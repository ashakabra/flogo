export module JsonSchema {

    export const 
        STRING = "string",
        NUMBER = "number",        
        ARRAY = "array";
        
    export class Types {

        public static toJsonType(type: string): any {
            let item = {};

            switch (type) {
                case "number":
                    item = JsonSchema.Types.numberType();
                    break;                
                case "string":
                case "option":
                case "user":
                case "date":
                case "datetime":
                    item = JsonSchema.Types.stringType();
                    break;
                case "array":
                    item = JsonSchema.Types.arrayType();
                    break;
                default:
                    console.log("Datatype case is not present :: "+type)
                    item = JsonSchema.Types.stringType();
            }
            return item;
        }
    
        public static numberType() {
            return {
                type: "number"
            };
        }
    
        public static stringType() {
            return {
                type: "string"
            };
        }


        public static arrayType() {
            return {
                "type": "array",
                "items": {
                  "type": "string"
                }
              };
        }        
    }
}