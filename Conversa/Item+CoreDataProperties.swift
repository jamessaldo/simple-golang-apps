//
//  Item+CoreDataProperties.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//
//

import Foundation
import CoreData


extension Item {

    @nonobjc public class func fetchRequest() -> NSFetchRequest<Item> {
        return NSFetchRequest<Item>(entityName: "Item")
    }

    @NSManaged public var content: String?
    @NSManaged public var timestamp: Date?
    @NSManaged public var title: String?

}

extension Item : Identifiable {

}
