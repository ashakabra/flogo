export module JsonSchema {

    export const 
        STRING = "string",
        NUMBER = "number",        
        ARRAY = "array";
        
    export class Types {

        public static toJsonType(schema: any, fieldName: string): any {
            let item = {};

            let jiraDataType = "";
            jiraDataType = schema.type;
            
            switch (jiraDataType) {
                case "number":
                    item = JsonSchema.Types.numberType(fieldName);
                    break;                
                case "string":
                case "option":
                case "user":
                case "date":
                case "datetime":
                    item = JsonSchema.Types.stringType(fieldName,jiraDataType);
                    break;
                case "array":
                    item = JsonSchema.Types.arrayType(fieldName,schema);
                    break;
                default:
                    console.log("Datatype case is not present :: "+jiraDataType)
                    item = JsonSchema.Types.stringType(fieldName, jiraDataType);
            }
            return item;
        }
    
        public static numberType(fieldName: string) {
            return {
                type: "number",
                flogoJiraID : fieldName 
            };
        }
    
        public static stringType(fieldName: string, jiraDatatype: string) {
            return {
                type: "string",
                flogoJiraID : fieldName,
                flogoJiraType : jiraDatatype
            };
        }


        public static arrayType(fieldName: string, schema: any) {
            let arrayItems = schema.items , arrayType = "";
            
            switch (arrayItems){
                case "version":
                    arrayType = "ArrayOfName";
                    break;
            }
            
            return {
                "type": "array",
                "items": {
                  "type": "string"
                },
                flogoJiraID : fieldName,
                flogoJiraType : arrayType
              };
        }        
    }
}