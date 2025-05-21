# Calling go from Snowpark

![SnowflakeGo](./snowflake-go.png)

# Overview 
An Example to mix go with python and use it in snowpark

Go takes an XML string, parses it, converts it to JSON, and returns it.
Python calls the Go function and directly receives a string an return a Python dict.


# Going to Go

The code for this example can be found in `xmltojson.go`

To build Go Shared Library run:
```
go build -o libxmljson.so -buildmode=c-shared xmltojson.go
```

And you’ll get:

`libxmljson.so` – the shared library

`libxmljson.h` – header file (you don't need this for Python)

> NOTE: CPU architecture is important. In my case I ran this code from Github Codespaces which
uses x84 architecture, so I know that is the the library architecture.

# Going Python

In Python we will Call Go, Get Dictionary

This can be done with some simple code that will mostly load the shared library, specify the arguments and voila! we are calling Go.

```
import ctypes
import json

# Load Go shared library
lib = ctypes.CDLL('./libxmljson.so')
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

```


# Going Snowpark

Ok fine. Well went to Go and Python, now lets go to Snowflake. For that we can mostly leverage our current python code. The biggest change is how we acceess the go lib. 
Mostly we provision this library in an snowflake stage. And we reference it into our Snowpark Python Function.

When adding imports to our UDF they are provision in an special `snowflake_import_directory` folder. So we adjust the code to read the `.so` from that location.

The other aspect to consider is to use `RESOURCE_CONSTRAINT=(architecture='x86')` we need that because this lib was compiled in x86, and the rest well just works.

I am so proud of the Snowflake Engineering team as they have found way to allow you to bring third party libraries and even binaries and run them [safely](https://medium.com/snowflake/snowpark-protection-through-java-scala-and-python-isolation-f8d10be61d56) in snowflake. 


```
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
```
