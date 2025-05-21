import ctypes
import json

import_folder = "./"
# Load Go shared library
lib = ctypes.CDLL(f'{import_folder}libxmljson.so')
lib.XMLToJSON.argtypes = [ctypes.c_char_p]
lib.XMLToJSON.restype = ctypes.c_char_p

# XML input
xml_input = """
<person>
    <name>Cristina</name>
    <age>30</age>
</person>
"""

# Call Go function
output_ptr = lib.XMLToJSON(xml_input.encode('utf-8'))

# Convert result to Python dict
output_str = ctypes.string_at(output_ptr).decode('utf-8')
parsed_dict = json.loads(output_str)

print(parsed_dict)
