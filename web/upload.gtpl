<html>
<head>
       <title>Add Command</title>

</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">

    <b>Only mp3 files are supported (no weird characters or whitespace)</b></br></br>
    Desired command: 
    <input type="text" name="command" /></br></br>
    <input type="file" name="uploadfile" /></br></br>
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="Add Command"/>
</form>


</body>
</html>