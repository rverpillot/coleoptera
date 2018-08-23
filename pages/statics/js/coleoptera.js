var map_individus;
var gmarkers = [];

Gp.Services.getConfig({
    apiKey: '1zt39dn13glty5q8zjcbcsbs'
});

function mapIGN(map) {
    L.geoportalLayer.WMTS({
        layer: "ORTHOIMAGERY.ORTHOPHOTOS"
    }).addTo(map);    
    L.geoportalLayer.WMTS({
        layer: "GEOGRAPHICALGRIDSYSTEMS.MAPS"
    }).addTo(map);

    var layerSwitcher = L.geoportalControl.LayerSwitcher();
    map.addControl(layerSwitcher);

    var mp = L.geoportalControl.MousePosition();
    map.addControl(mp);
}

function createMap(tag, center, zoom) {
    if ($(tag).length == 0) return;

    map_individus = L.map($(tag)[0]).setView([center.lat, center.lng], zoom);
    mapIGN(map_individus)

    var center = map_individus.getCenter()
    var data = { lat: center.lat, lng: center.lng, zoom: map_individus.getZoom() }
    sendMsg(null, "map-loaded", "map", "page", data)

    console.log("createMap")
}

function refreshMap(center, zoom) {
    // google.maps.event.trigger(map_individus, 'resize')
    map_individus.setView(center, zoom)
}

function showMarkers(tag, markers) {

    console.log("show markers")

    $.each(gmarkers, function (i, marker) {
        marker.remove()
    })

    gmarkers = []
    $.each(markers, function (i, location) {
        var position = location.Location
        // if (!map_individus.getBounds().contains(position)) {
        //     return
        // }
        // var nb = location.Infos.length
        // var label = nb.toString()
        // if (nb > 9) {
        //     label = "#"
        // }
        var marker = L.marker(position, {
            title: location.Infos.join("\n")
        })
        gmarkers.push(marker)
        marker.addTo(map_individus)
    })
}

var editMap;
var editMarker;
function updateEditMap(latitude, longitude) {
    var pos = { lat: latitude, lng: longitude }
    editMarker.setPosition(pos)
    editMap.setCenter(pos)
}

function createEditMap(tag) {
    var latitude = parseFloat($("[name=latitude]").val())
    var longitude = parseFloat($("[name=longitude]").val())
    var position = { lat: 47.626951, lng: 6.997541 }

    if (!isNaN(longitude) && !isNaN(latitude)) {
        position = { lat: latitude, lng: longitude };
    }
    editMap = L.map($(tag)[0], {
        scrollWheelZoom: false
    }).setView(position, 12);
    mapIGN(editMap);

    editMarker = L.marker(position, {
        // title: position.toString(),
        draggable: true
    });
    editMarker.addTo(editMap);
    console.log("createEditMap")

    editMarker.on('dragend', function (event) {
        var position = editMarker.getLatLng()
        // editMarker.setTitle(position.toString())
        sendMsg(null, "position", "map", "page", position)
    })
}

function createPreviewMap(tag, longitude, latitude) {
    var position = { lat: latitude, lng: longitude };

    var map = L.map($(tag)[0], {
        scrollWheelZoom: false 
    }).setView(position, 8);
    mapIGN(map);

    var marker = L.marker(position, {
        title: position.toString()
    })
    marker.addTo(map)
    console.log("createPreviewMap")
}


$(document).on("ihui:display", function (ev, tag) {
    //    console.log(tag)

    $(tag).find('.ui.modal').modal({ closable: false }).modal("show")
    $(tag).find('.ui.checkbox').checkbox()
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
    $(tag).find('.ui.dropdown.additions').dropdown({
        forceSelection: false,
        allowAdditions: true,
        hideAdditions: false,
        message: {
            addResult: 'Ajouter <b>{term}</b>',
            count: '{count} selectionné(s)',
            maxSelections: 'Max {maxCount} selections',
            noResults: 'Aucun résultat.'
        }
    })

    $(tag).find('select.dropdown').dropdown()

})
