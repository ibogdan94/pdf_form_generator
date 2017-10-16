// var stage = new Kinetic.Stage({
//     container: 'container',
//     width: 300,
//     height: 300
// });
//
// var layer = new Kinetic.Layer();
// stage.add(layer);
//
// var image = new Image();
// image.src = "https://pp.userapi.com/c639619/v639619579/43af0/wJ3cOZ7WQVg.jpg";
// image.setAttribute('crossOrigin', 'anonymous');
// // image.crossOrigin = 'anonymous';
//
// $('#imgSource').on('change', function () {
//     image.src = $(this).val();
// });
//
// image.onload = function () {
//     var image1 = new Kinetic.Image({
//         x: 0,
//         y: 0,
//         width: 300,
//         height: 300,
//         image: image
//     });
//     layer.add(image1);
//     layer.draw();
// };
//
// var currentLabel;
//
// $("#addbutton").click(function () {
//     // simple label
//     var label = new Kinetic.Label({
//         x: 20,
//         y: 20,
//         draggable: true
//     });
//
//     label.add(new Kinetic.Tag({
//         fill: 'green'
//     }));
//
//     label.add(new Kinetic.Text({
//         text: $("#newtext").val(),
//         fontFamily: 'Verdana',
//         fontSize: 18,
//         padding: 10,
//         fill: 'white'
//     }));
//
//     layer.add(label);
//
//     layer.draw();
//
//     label.on('mousedown', function (e) {
//         console.log(e);
//     });
//
//     label.on('mouseup', function (e) {
//         console.log(e);
//         label.className = 'active';
//         $("#newtext").val(e.target.attrs.text);
//         currentLabel = label;
//     });
//
//     $("#newtext").val('');
// });
//
// $('#update_prop').on('click', function (e) {
//     console.log(currentLabel);
//     // currentLabel.attrs.text = $("#newtext").val();
//     currentLabel.attrs.fontSize = parseInt($("#font_size").val());
//     // currentLabel.target.attr.fill = $("#font_color").val();
// });
//
// $('#save').on('click', function () {
//     // stage.toDataURL({
//     //     callback: function(dataUrl) {
//     //         window.open(dataUrl);
//     //     }
//     // });
//
//     console.log(layer.getCanvas().toDataURL());
//
//     $.ajax({
//         type: "POST",
//         url: "/save",
//         data: {"b64": layer.getCanvas().toDataURL() }
//     }).done(function(o) {
//         console.log('saved');
//     });
// });


(function() {
    var $ = function(id){return document.getElementById(id)};

    document.querySelectorAll("canvas").forEach(function(canvas) {

        // canvas.addEventListener('click', function(event) {
        //     console.log(123);
        //     this.className += " active";
        // }, false);

        var CANVAS = this.__canvas = new fabric.Canvas(canvas);


        CANVAS.setBackgroundImage(
            // 'https://pp.userapi.com/c639619/v639619579/43af0/wJ3cOZ7WQVg.jpg',
            canvas.getAttribute("data-src"),
            CANVAS.renderAll.bind(CANVAS),
            {
                // height: '297mm',
                // width: '210mm',
                // opacity: 0.5,
                // angle: 45,
                // left: 400,
                // top: 400,
                // originX: 'left',
                // originY: 'top',
                crossOrigin: 'anonymous'
            }
        );

        fabric.Object.prototype.transparentCorners = false;

        var addTextButton = $('addbutton');

        addTextButton.onclick = function () {
            var newText = $('newtext').value;

            var text = new fabric.Text(newText, {
                fontSize: parseInt($('fontSize-control').value, 10),
                fontWeight: parseInt($('fontWeight-control').value, 10),
                charSpacing: parseInt($('charSpacing-control').value, 10)
            });

            text.setColor('#' + $('color-control').value);

            CANVAS.add(text);
            testAAA(text);

            $('newtext').value = "";
        };

        $("delete").onclick = function (event) {
            var activeObject = CANVAS.getActiveObject(),
                activeGroup = CANVAS.getActiveGroup();
            if (activeObject) {
                if (confirm('Are you sure?')) {
                    CANVAS.remove(activeObject);
                }
            } else if (activeGroup) {
                if (confirm('Are you sure?')) {
                    var objectsInGroup = activeGroup.getObjects();
                    CANVAS.discardActiveGroup();
                    objectsInGroup.forEach(function(object) {
                        CANVAS.remove(object);
                    });
                }
            }
        };

        var testAAA = function (elem) {
            var angleControl = $('angle-control');
            angleControl.oninput = function() {
                elem.set('angle', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var scaleControl = $('scale-control');
            scaleControl.oninput = function() {
                elem.scale(parseFloat(this.value)).setCoords();
                fontSizeControl.value = (elem.fontSize * elem.scaleX).toFixed(0);
                CANVAS.renderAll();
            };

            var topControl = $('top-control');
            topControl.oninput = function() {
                elem.set('top', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var leftControl = $('left-control');
            leftControl.oninput = function() {
                elem.set('left', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var skewXControl = $('skewX-control');
            skewXControl.oninput = function() {
                elem.set('skewX', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var skewYControl = $('skewY-control');
            skewYControl.oninput = function() {
                elem.set('skewY', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var charSpacingControl = $('charSpacing-control');
            charSpacingControl.oninput = function() {
                elem.set('charSpacing', parseInt(this.value, 10)).setCoords();
                CANVAS.renderAll();
            };

            var fontSizeControl = $('fontSize-control');
            fontSizeControl.onchange = function() {
                elem.set('fontSize', parseFloat(this.value));
                CANVAS.renderAll();
            };

            var colorControl = $('color-control');
            colorControl.onchange = function() {
                elem.set('fill', '#' + this.value);
                CANVAS.renderAll();
            };

            var fontWeightControl = $('fontWeight-control');
            fontWeightControl.onchange = function() {
                console.log(elem, this.value);
                elem.set('fontWeight', parseInt(this.value, 10));
                CANVAS.renderAll();
            };

            function updateControls() {
                scaleControl.value = elem.scaleX;
                angleControl.value = elem.angle;
                leftControl.value = elem.left;
                topControl.value = elem.top;
                skewXControl.value = elem.skewX;
                skewYControl.value = elem.skewY;
                charSpacingControl.value = elem.charSpacing;
                colorControl.value = elem.fill.substring(1);
                fontSizeControl.value = elem.fontSize;
                fontWeightControl.value = elem.fontWeight;
            }

            CANVAS.on({
                'object:moving': updateControls,
                'object:scaling': updateControls,
                'object:resizing': updateControls,
                'object:rotating': updateControls,
                'object:skewing': updateControls,
                'object:modified': updateControls
            });
        };

        var sendAjax = function (method, url, body) {
            var xmlhttp = new XMLHttpRequest();

            xmlhttp.onreadystatechange = function() {
                if (xmlhttp.readyState === XMLHttpRequest.DONE ) {
                    if (xmlhttp.status >= 200 && xmlhttp.status < 300) {
                        const data = JSON.parse(xmlhttp.response);
                        if (data.url) {
                            var win = window.open(data.url, '_blank');
                        }
                    } else if (xmlhttp.status === 400) {
                        alert('There was an error 400');
                    } else {
                        alert('something else other than 200 was returned');
                    }
                }
            };

            xmlhttp.open(method, url, true);
            xmlhttp.send(JSON.stringify(body));
        };


        $('save').onclick = function () {
            sendAjax("POST", "/pdf/save", {"b64": CANVAS.toDataURL()});
        };
    });




    // var canvas = new fabric.Canvas('canvas');
    // $(document).ready(function() {
    //     $('#text-font-size').keyup(function() {
    //         var val = $(this).val();
    //         if (isNaN(val)) {
    //             alert('please enter number');
    //             $(this).val('');
    //         }
    //         var activeObject = canvas.getActiveObject();
    //         activeObject.fontSize = val;
    //         canvas.renderAll();
    //     });
    //     $('#add-text-btn').click(function() {
    //         if ($('#text-font-size').val()) {
    //             var txtfontsize = $('#text-font-size').val();
    //         } else {
    //             var txtfontsize = 40;
    //         }
    //         var message = $('#add-text-value').val();
    //         //var txtfontfamily = $('#font-family').val();
    //         var new_text = new fabric.IText(message, {
    //             left: 100,
    //             top: 100,
    //             fontSize: txtfontsize,
    //             lockUniScaling: true,
    //             //fontFamily: txtfontfamily,
    //             fill: '#000'
    //         });
    //         canvas.add(new_text);
    //         canvas.setActiveObject(new_text);
    //     });
    //
    //     canvas.on('object:selected', function(options) {
    //         if (options.target) {
    //             $("textarea#add-text-value").val(options.target.text);
    //             $("#text-font-size").val(options.target.fontSize);
    //         }
    //     });
    //
    //     canvas.on('object:scaling', function(event) {
    //         if (event.target) {
    //             $("textarea#add-text-value").val(event.target.text);
    //             $("#text-font-size").val((event.target.fontSize * event.target.scaleX).toFixed(0));
    //         }
    //     });
    //
    //     canvas.on('object:modified', function(event) {
    //         if (event.target) {
    //             event.target.fontSize *= event.target.scaleX;
    //             event.target.fontSize = event.target.fontSize.toFixed(0);
    //             event.target.scaleX = 1;
    //             event.target.scaleY = 1;
    //             event.target._clearCache();
    //             $("textarea#add-text-value").val(event.target.text);
    //             $("#text-font-size").val(event.target.fontSize);
    //         }
    //     });
    //
    // });

    // $('#save').on('click', function () {
    //     // stage.toDataURL({
    //     //     callback: function(dataUrl) {
    //     //         window.open(dataUrl);
    //     //     }
    //     // });
    //
    //     console.log(layer.getCanvas().toDataURL());
    //
    //     $.ajax({
    //         type: "POST",
    //         url: "/save",
    //         data: {"b64": layer.getCanvas().toDataURL() }
    //     }).done(function(o) {
    //         console.log('saved');
    //     });
    // });
})();