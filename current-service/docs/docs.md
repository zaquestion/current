# current

## current.proto

### Messages

<a name="PostLocationBigBrotherRequest"></a>

#### PostLocationBigBrotherRequest

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| latitude | TYPE_DOUBLE | 1 |  |
| longitude | TYPE_DOUBLE | 2 |  |
| accuracy | TYPE_DOUBLE | 3 |  |
| altitude | TYPE_DOUBLE | 4 |  |
| bearing | TYPE_DOUBLE | 5 |  |
| speed | TYPE_DOUBLE | 6 |  |
| battlevel | TYPE_INT32 | 7 |  |
| time | TYPE_STRING | 8 |  |

<a name="PostLocationTaskerRequest"></a>

#### PostLocationTaskerRequest

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| location | TYPE_DOUBLE | 1 |  |
| speed | TYPE_DOUBLE | 2 |  |
| battery | TYPE_INT32 | 3 |  |
| charging | TYPE_BOOL | 4 |  |
| time | TYPE_STRING | 5 |  |

<a name="Error"></a>

#### Error

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| err | TYPE_STRING | 1 |  |

<a name="GetLocationRequest"></a>

#### GetLocationRequest


<a name="Location"></a>

#### Location

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| latitude | TYPE_DOUBLE | 1 |  |
| longitude | TYPE_DOUBLE | 2 |  |
| speed | TYPE_DOUBLE | 3 |  |
| battery | TYPE_INT32 | 4 |  |
| charging | TYPE_BOOL | 5 |  |
| last_updated | TYPE_STRING | 6 |  |
| err | TYPE_STRING | 7 |  |

### Services

#### Current

| Method Name | Request Type | Response Type | Description|
| ---- | ---- | ------------ | -----------|
| PostLocationBigBrother | PostLocationBigBrotherRequest | Error |  |
| PostLocationTasker | PostLocationTaskerRequest | Error |  |
| GetLocation | GetLocationRequest | Location |  |

#### Current - Http Methods

##### POST `/location/bigbrother`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |
| latitude | query | TYPE_DOUBLE |
| longitude | query | TYPE_DOUBLE |
| accuracy | query | TYPE_DOUBLE |
| altitude | query | TYPE_DOUBLE |
| bearing | query | TYPE_DOUBLE |
| speed | query | TYPE_DOUBLE |
| battlevel | query | TYPE_INT32 |
| time | query | TYPE_STRING |

##### POST `/location/tasker`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |
| location | body | TYPE_DOUBLE |
| speed | body | TYPE_DOUBLE |
| battery | body | TYPE_INT32 |
| charging | body | TYPE_BOOL |
| time | body | TYPE_STRING |

##### GET `/location`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |


<style type="text/css">

body{
    font-family      : helvetica, arial, freesans, clean, sans-serif;
    color            : #003269;
    background-color : #fff;
    border-color     : #999999;
    border-width     : 2px;
    line-height      : 1.5;
    margin           : 2em 3em;
    text-align       :left;
    font-size        : 16px;
    padding          : 0 100px 0 100px;

    width         : 1024px;
    margin-top    : 0px;
    margin-bottom : 2em;
    margin-left   : auto;
    margin-right  : auto;
}

h1 {
    font-family : 'Gill Sans Bold', 'Optima Bold', Arial, sans-serif;
    color       : #577AD3;
    font-weight : 400;
    font-size   : 48px;
}
h2{
    margin-bottom : 1em;
    padding-top   : 0.5em;
    color         : #003269;
    font-size     : 36px;
}
h3{
    border-bottom : 1px dotted #aaa;
    color         : #4660A4;
    font-size     : 30px;
}
h4 {
    font-size: 24px;
}
h5 {
    font-size: 18px;
}
code {
    font-family      : Consolas, "Inconsolata", Menlo, Monaco, Lucida Console, Liberation Mono, DejaVu Sans Mono, Bitstream Vera Sans Mono, Courier New, monospace, serif; /* Taken from the stackOverflow CSS*/
    background-color : #f5f5f5;
    border           : 1px solid #e1e1e8;
}


pre {
    display          : block;
    background-color : #f5f5f5;
    border           : 1px solid #ccc;
    padding          : 3px 3px 3px 3px;
}
pre code {
    white-space      : pre-wrap;
    padding          : 0;
    border           : 0;
    background-color : code;
}

table {
	border-collapse: collapse; border-spacing: 0;
	width: 100%;
	margin-bottom : 3em;
}
td, th {
	vertical-align: top;
	padding: 4px 10px;
	border: 1px solid #9BC3EB;
}
tr:nth-child(even) td, tr:nth-child(even) th {
	background: #EBF4FE;
}
th:nth-child(4) {
	width: auto;
}

</style>
