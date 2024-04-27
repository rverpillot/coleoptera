
$(document).on("element-created element-updated", function (event) {
    // console.log("element-created element-updated")
    doSemanticUI($("#" + event.detail.element))
})


var map_individus;

function showMarkers(tag, markers) {
    // console.log("show markers")
    var gmarkers = []
    $.each(markers, function (i, marker) {
        var marker = {
            position: { x: marker.Location.lng, y: marker.Location.lat, projection: "CRS:84" },
            content: marker.Infos.join("\n"),
            url: "images/icon.png",
        }
        gmarkers.push(marker)
    })

    map_individus.setMarkersOptions(gmarkers)
}

function createMap(tag, center, zoom, pageName) {
    if ($(tag).length == 0) return;

    map_individus = Gp.Map.load(
        $(tag)[0],
        {
            apiKey: "cartes,satellite,ortho",
            center: { geolocate: true },
            zoom: zoom,
            mapEventsOptions: {
            },
            controlsOptions: {
                "search": {
                    maximised: true
                },
                "layerswitcher": {
                    maximised: true
                },
            },
            layersOptions: {
                "ORTHOIMAGERY.ORTHOPHOTOS": {},
                "GEOGRAPHICALGRIDSYSTEMS.PLANIGNV2": {},
                "GEOGRAPHICALGRIDSYSTEMS.MAPS": {},
            }
        })

    // map_individus.on("moveend zoomend", function (ev) {
    //     var center = map_individus.getCenter()
    //     var data = { lat: center.lat, lng: center.lng, zoom: map_individus.getZoom() }
    //     ihui.trigger("map-changed", pageName, data, false)
    // })

    console.log("createMap")
}

function refreshMap(center, zoom) {
    console.log("refreshMap")
    //    map_individus.setView(center, zoom)
    map_individus.invalidateSize()
}


function createEditMap(tag, pageName) {
    var latitude = parseFloat($("[name=latitude]").val())
    var longitude = parseFloat($("[name=longitude]").val())
    var position = { y: 47.626951, x: 6.997541, projection: "CRS:84"}

    if (!isNaN(longitude) && !isNaN(latitude)) {
        position.x = longitude;
        position.y = latitude;
    }

    var editMap = Gp.Map.load(
        $(tag)[0],
        {
            apiKey: "cartes,satellite,ortho",
            center: position,
            zoom: 10,
            mapEventsOptions: {
            },
            controlsOptions: {
                "search": {
                    maximised: true
                },
                "layerswitcher": {
                    maximised: true
                },
            },
            layersOptions: {
                "ORTHOIMAGERY.ORTHOPHOTOS": {},
                "GEOGRAPHICALGRIDSYSTEMS.PLANIGNV2": {},
                "GEOGRAPHICALGRIDSYSTEMS.MAPS": {},
            },
            markerOptions: [position]
        })

    var editMarker = L.marker(position, {
        // title: position.toString(),
        draggable: true
    });
    editMarker.addTo(editMap);
    // console.log("createEditMap")

    editMarker.on('dragend', function (event) {
        var position = editMarker.getLatLng()
        // editMarker.setTitle(position.toString())
        ihui.trigger("position", pageName, position)
    })

    editMap.on("moveend", function (ev) {
        var pos = editMap.getCenter()
        editMarker.setLatLng(pos)
    })

}

function createPreviewMap(tag, longitude, latitude) {
    var position = { lat: latitude, lng: longitude };

    var previewMap = L.map($(tag)[0], {
        center: position,
        zoom: 6,
        scrollWheelZoom: false
    })
    mapIGN(previewMap, false);

    var marker = L.marker(position, {})
    marker.addTo(previewMap)
    // console.log("createPreviewMap")
}


function doSemanticUI(tag) {
    $(tag).find('.ui.modal').modal({ closable: false }).modal("show")
    $(tag).find('.ui.dropdown').dropdown({
        forceSelection: false,
        fullTextSearch: true,
        message: {
            addResult: "Ajouter <b>{term}</b>",
            count: "{count} selectionné(s)",
            maxSelections: 'Max {maxCount} selections',
            noResults: 'Aucun résultat.'
        }
    })
    $(tag).find('select.dropdown').dropdown()

}
