
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
# current

## current.proto

### Messages

<a name="PostLocationFromBigBrotherRequest"></a>

#### PostLocationFromBigBrotherRequest

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| latitude | TYPE_DOUBLE | 1 |  |
| longitude | TYPE_DOUBLE | 2 |  |
| accuracy | TYPE_DOUBLE | 3 |  |
| altitude | TYPE_DOUBLE | 4 |  |
| bearing | TYPE_DOUBLE | 5 |  |
| speed | TYPE_DOUBLE | 6 |  |

<a name="PostLocationFromBigBrotherReply"></a>

#### PostLocationFromBigBrotherReply

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| err | TYPE_STRING | 2 |  |

<a name="GetLocationRequest"></a>

#### GetLocationRequest


<a name="GetLocationReply"></a>

#### GetLocationReply

| Name | Type | Field Number | Description|
| ---- | ---- | ------------ | -----------|
| latitude | TYPE_DOUBLE | 1 |  |
| longitude | TYPE_DOUBLE | 2 |  |
| err | TYPE_STRING | 3 |  |

### Services

                Name: CurrentService
                Desc: 
| Method Name | Request Type | Response Type | Description|
| ---- | ---- | ------------ | -----------|
| PostLocationFromBigBrother | PostLocationFromBigBrotherRequest | PostLocationFromBigBrotherReply |  |
| GetLocation | GetLocationRequest | GetLocationReply |  |

#### CurrentService - Http Methods

##### GET `/location/bigbrother`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |
| latitude | query | TYPE_DOUBLE |
| longitude | query | TYPE_DOUBLE |
| accuracy | query | TYPE_DOUBLE |
| altitude | query | TYPE_DOUBLE |
| bearing | query | TYPE_DOUBLE |
| speed | query | TYPE_DOUBLE |

##### GET `/location`



| Parameter Name | Location | Type |
| ---- | ---- | ------------ |
