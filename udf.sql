CREATE OR REPLACE FUNCTION XML_TO_JSON(XMLINPUT VARCHAR) 
RETURNS OBJECT
LANGUAGE PYTHON
RESOURCE_CONSTRAINT=(architecture='x86')
RUNTIME_VERSION = 3.11
IMPORTS = ('@mystage/libxmljson.so')
HANDLER = 'main'
AS 
$$
import ctypes
import json
import sys
def main(xmlinput):
    IMPORT_DIRECTORY_NAME = "snowflake_import_directory"
    import_folder = sys._xoptions[IMPORT_DIRECTORY_NAME]

    # Load Go shared library
    lib = ctypes.CDLL(f'{import_folder}libxmljson.so')
    lib.XMLToJSON.argtypes = [ctypes.c_char_p]
    lib.XMLToJSON.restype = ctypes.c_char_p

    output_ptr = lib.XMLToJSON(xmlinput.encode('utf-8'))

    # Convert result to Python dict
    output_str = ctypes.string_at(output_ptr).decode('utf-8')
    return json.loads(output_str)
$$;


select  XML_TO_JSON($$
<person>
    <name>Cristina</name>
    <age>30</age>
</person>
$$);